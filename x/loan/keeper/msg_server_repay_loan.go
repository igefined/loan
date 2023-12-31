package keeper

import (
	"context"

	"loan/x/loan/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) RepayLoan(goCtx context.Context, msg *types.MsgRepayLoan) (*types.MsgRepayLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	loan, ok := k.GetLoan(ctx, msg.Id)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key %d not found", msg.Id)
	}

	if loan.State != "approved" {
		return nil, sdkerrors.Wrapf(types.ErrWrongLoanState, "%v", loan.State)
	}

	lender, _ := sdk.AccAddressFromBech32(loan.Lender)
	borrower, _ := sdk.AccAddressFromBech32(loan.Borrower)
	if msg.Creator != loan.Borrower {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "cannot repay: not the borrower")
	}

	amount, _ := sdk.ParseCoinsNormalized(loan.Amount)
	fee, _ := sdk.ParseCoinsNormalized(loan.Fee)
	collateral, _ := sdk.ParseCoinsNormalized(loan.Collateral)

	if err := k.bankKeeper.SendCoins(ctx, borrower, lender, amount); err != nil {
		return nil, err
	}

	if err := k.bankKeeper.SendCoins(ctx, borrower, lender, fee); err != nil {
		return nil, err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, borrower, collateral); err != nil {
		return nil, err
	}

	loan.State = "repayed"
	k.SetLoan(ctx, loan)

	return &types.MsgRepayLoanResponse{}, nil
}
