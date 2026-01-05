// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "./InsightToken.sol";
import "./lib/SafeMath.sol";

contract RewardPool is AccessControl, Pausable, ReentrancyGuard {
    using SafeMath for uint256;

    byte32 public constant DISTRIBUTOR_ROLE = keccak256("DISTRIBUTOR_ROLE");
    byte32 public constant OPERATR_ROLE = keccak256("OPERATR_ROLE");

    InsightToken public immutable insightToken;
    IERC20 public immutable usdcToken;

    // 黑洞地址 ：销毁代币用
    address public constant BURN_ADDRESS =
        0x000000000000000000000000000000000000dEaD;

    // 分配比例,基数  10000
    uint256 public buybackRate = 4000; // 40% 回购比例
    uint256 public rewardRate = 3000; // 30% 奖励比例
    uint256 public treasuryRate = 2000; // 20% 入国库比例
    uint256 public burnRate = 1000; // 10% 销毁比例，减少总供应，提升币价

    // 记录待分配的利润
    struct PendingProfit {
        uint256 usdcAmount;
        uint256 tokenAmount;
    }
    PendingProfit public pendingProfit;

    // 获利 事件
    event ProfitReceived(uint usdcAmount, uint256 tokenAmount);
    // 回购事件
    event BuyBackExcuted(
        uint usdcSpent,
        uint256 tokenBougnt,
        uint256 tokenBurned
    );
    // 奖励事件
    event RewardDistribution(address[] receipts, uint256 amount);
    // 配置更新事件
    event AllocationUpdated(uint256 amount);

    constructor(address _insightToken, address _usdcToken, address _treasury) {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(OPERATR_ROLE, msg.sender);
        _grantRole(DISTRIBUTOR_ROLE, msg.sender);

        insightToken = InsightToken(_insightToken);
        usdcToken = IERC20(_usdcToken);
        treasury = _treasury;
    }

    /**
     * 后端调用此方法，获取收益
     * @param usdcAmount USDC数量
     * @param tokenAmount 平台代币数量
     */
    function receiveProfit(
        uint256 usdcAmount,
        uint256 tokenAmount
    ) external onlyRole(DISTRIBUTOR_ROLE) whenNotPaused {
        if (usdcAmount > 0) {
            require(
                usdcToken.transferFrom(msg.sender, treasury, usdcAmount),
                "transfer usdc failed"
            );
            pendingProfit.usdcAmount += usdcAmount;
        }
        if (tokenAmount > 0) {
            require(
                insightToken.transferFrom(msg.sender, treasury, tokenAmount),
                "transfer token failed"
            );
            pendingProfit.tokenAmount += tokenAmount;
        }

        emit ReceiveProfit(usdcAmount, tokenAmount);
    }

    function RewardDistributionistribution(
        address[] calldata receipts,
        uint256[] calldata amounts
    ) external onlyRole(DISTRIBUTOR_ROLE) whenNotPaused nonReentrant {
        require(
            receipts.length == amounts.length,
            "receipts and amounts length mismatch"
        );

        // 获取待总的待分配的利润
        uint256 totalUsdc = pendingProfit.usdcAmount;
        uint256 totalToken = pendingProfit.tokenAmount;

        // 1、按比例回购（销毁）平台代币
        uint256 buyBackAmount = totalToken.mul(buybackRate).div(10000);
        if (buyBackAmount > 0) {
            // 实际生产环境应该桶 dex 执行链上回购，拿到回购代币再销毁
            // 这里采用直接销毁合约中已有的部分代币
            insightToken.transfer(BURN_ADDRESS, buyBackAmount);
            pendingProfit.totalToken -= buyBackAmount; // 更新待分配余额
        }

        // 2、分发代币奖励给用户
        uint256 rewardAmount = totalToken.mul(rewardRate).div(10000);

        for (uint256 i = 0; i < receipts.length; i++) {
            insightToken.transferFrom(msg.sender, receipts[i], amounts[i]);
        }
    }

    /**
     * 暂停功能: 暂停合约的所有可暂停功能
     */
    function pause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _pause();
    }

    function unpause() external onlyRole(DEFAULT_ADMIN_ROLE) {
        _unpause();
    }
}
