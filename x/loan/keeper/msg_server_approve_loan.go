package keeper

import (
	"context"
	"loan/x/loan/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) ApproveLoan(goCtx context.Context, msg *types.MsgApproveLoan) (*types.MsgApproveLoanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	loan, ok := k.GetLoan(ctx, msg.Id)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "key %d does not exist", msg.Id)
	}

	if loan.State != "requested" {
		return nil, sdkerrors.Wrapf(types.ErrWrongLoanState, "%v", loan)
	}

	lender, _ := sdk.AccAddressFromBech32(msg.Creator)
	borrower, _ := sdk.AccAddressFromBech32(loan.Borrower)
	amount, err := sdk.ParseCoinsNormalized(loan.Amount)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrWrongLoanState, "cannot parse coins in loan amount")
	}

	err = k.bankKeeper.SendCoins(ctx, lender, borrower, amount)
	if err != nil {
		return nil, err
	}

	loan.Lender = msg.Creator
	loan.State = "approved"
	k.SetLoan(ctx, loan)

	return &types.MsgApproveLoanResponse{}, nil
}
