package keeper_test

import (
	"cosmossdk.io/math"
	i "github.com/KYVENetwork/chain/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// Pool
	"github.com/KYVENetwork/chain/x/pool/types"
	// Gov
	govV1Types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

/*

TEST CASES - msg_server_update_params.go

* Check default params
* Invalid authority (transaction)
* Invalid authority (proposal)
* Update every param at once
* Update no param
* Update with invalid formatted payload

* Update protocol inflation share
* Update protocol inflation share with invalid value

* Update pool inflation payout rate
* Update pool inflation payout rate with invalid value

* Update max voting power per pool
* Update max voting power per pool with invalid value

*/

var _ = Describe("msg_server_update_params.go", Ordered, func() {
	s := i.NewCleanChain()

	gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()

	params, _ := s.App().GovKeeper.Params.Get(s.Ctx())
	minDeposit := params.MinDeposit
	votingPeriod := params.VotingPeriod

	delegations, _ := s.App().StakingKeeper.GetAllDelegations(s.Ctx())
	voter := sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress)

	BeforeEach(func() {
		s = i.NewCleanChain()

		delegations, _ := s.App().StakingKeeper.GetAllDelegations(s.Ctx())
		voter = sdk.MustAccAddressFromBech32(delegations[0].DelegatorAddress)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Check default params", func() {
		// ASSERT
		params := s.App().PoolKeeper.GetParams(s.Ctx())

		Expect(params.ProtocolInflationShare).To(Equal(types.DefaultProtocolInflationShare))
		Expect(params.PoolInflationPayoutRate).To(Equal(types.DefaultPoolInflationPayoutRate))
		Expect(params.MaxVotingPowerPerPool).To(Equal(types.DefaultMaxVotingPowerPerPool))
	})

	It("Invalid authority (transaction)", func() {
		// ARRANGE
		msg := &types.MsgUpdateParams{
			Authority: i.DUMMY[0],
			Payload:   "{}",
		}

		// ACT
		_, err := s.RunTx(msg)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Invalid authority (proposal)", func() {
		// ARRANGE
		msg := &types.MsgUpdateParams{
			Authority: i.DUMMY[0],
			Payload:   "{}",
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		// ACT
		_, err := s.RunTx(proposal)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Update every param at once", func() {
		// ARRANGE
		payload := `{
			"protocol_inflation_share": "0.2",
			"pool_inflation_payout_rate": "0.05"
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		vote := govV1Types.NewMsgVote(
			voter, 1, govV1Types.VoteOption_VOTE_OPTION_YES, "",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)
		_, voteErr := s.RunTx(vote)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().PoolKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.ProtocolInflationShare).To(Equal(math.LegacyMustNewDecFromStr("0.2")))
		Expect(updatedParams.PoolInflationPayoutRate).To(Equal(math.LegacyMustNewDecFromStr("0.05")))
		Expect(updatedParams.MaxVotingPowerPerPool).To(Equal(types.DefaultMaxVotingPowerPerPool))
	})

	It("Update no params", func() {
		// ARRANGE
		payload := `{}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		vote := govV1Types.NewMsgVote(
			voter, 1, govV1Types.VoteOption_VOTE_OPTION_YES, "",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)
		_, voteErr := s.RunTx(vote)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().PoolKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.ProtocolInflationShare).To(Equal(types.DefaultProtocolInflationShare))
		Expect(updatedParams.PoolInflationPayoutRate).To(Equal(types.DefaultPoolInflationPayoutRate))
		Expect(updatedParams.MaxVotingPowerPerPool).To(Equal(types.DefaultMaxVotingPowerPerPool))
	})

	It("Update with invalid formatted payload", func() {
		// ARRANGE
		payload := `{
			"protocol_inflation_share: 20
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().PoolKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.ProtocolInflationShare).To(Equal(types.DefaultProtocolInflationShare))
		Expect(updatedParams.PoolInflationPayoutRate).To(Equal(types.DefaultPoolInflationPayoutRate))
		Expect(updatedParams.MaxVotingPowerPerPool).To(Equal(types.DefaultMaxVotingPowerPerPool))
	})

	It("Update protocol inflation share", func() {
		// ARRANGE
		payload := `{
			"protocol_inflation_share": "0.07"
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		vote := govV1Types.NewMsgVote(
			voter, 1, govV1Types.VoteOption_VOTE_OPTION_YES, "",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)
		_, voteErr := s.RunTx(vote)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().PoolKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.ProtocolInflationShare).To(Equal(math.LegacyMustNewDecFromStr("0.07")))
		Expect(updatedParams.PoolInflationPayoutRate).To(Equal(types.DefaultPoolInflationPayoutRate))
		Expect(updatedParams.MaxVotingPowerPerPool).To(Equal(types.DefaultMaxVotingPowerPerPool))
	})

	It("Update protocol inflation share with invalid value", func() {
		// ARRANGE
		payload := `{
			"protocol_inflation_share": "invalid"
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().PoolKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.ProtocolInflationShare).To(Equal(types.DefaultProtocolInflationShare))
		Expect(updatedParams.PoolInflationPayoutRate).To(Equal(types.DefaultPoolInflationPayoutRate))
		Expect(updatedParams.MaxVotingPowerPerPool).To(Equal(types.DefaultMaxVotingPowerPerPool))
	})

	It("Update pool inflation payout rate", func() {
		// ARRANGE
		payload := `{
			"pool_inflation_payout_rate": "0.2"
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		vote := govV1Types.NewMsgVote(
			voter, 1, govV1Types.VoteOption_VOTE_OPTION_YES, "",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)
		_, voteErr := s.RunTx(vote)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().PoolKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.ProtocolInflationShare).To(Equal(types.DefaultProtocolInflationShare))
		Expect(updatedParams.PoolInflationPayoutRate).To(Equal(math.LegacyMustNewDecFromStr("0.2")))
		Expect(updatedParams.MaxVotingPowerPerPool).To(Equal(types.DefaultMaxVotingPowerPerPool))
	})

	It("Update pool inflation payout rate with invalid value", func() {
		// ARRANGE
		payload := `{
			"pool_inflation_payout_rate": "1.2"
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().PoolKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.ProtocolInflationShare).To(Equal(types.DefaultProtocolInflationShare))
		Expect(updatedParams.PoolInflationPayoutRate).To(Equal(types.DefaultPoolInflationPayoutRate))
		Expect(updatedParams.MaxVotingPowerPerPool).To(Equal(types.DefaultMaxVotingPowerPerPool))
	})

	It("Update max voting power per pool", func() {
		// ARRANGE
		payload := `{
			"max_voting_power_per_pool": "0.2"
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		vote := govV1Types.NewMsgVote(
			voter, 1, govV1Types.VoteOption_VOTE_OPTION_YES, "",
		)

		// ACT
		_, submitErr := s.RunTx(proposal)
		_, voteErr := s.RunTx(vote)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().PoolKeeper.GetParams(s.Ctx())

		Expect(submitErr).NotTo(HaveOccurred())
		Expect(voteErr).NotTo(HaveOccurred())

		Expect(updatedParams.ProtocolInflationShare).To(Equal(types.DefaultProtocolInflationShare))
		Expect(updatedParams.PoolInflationPayoutRate).To(Equal(types.DefaultPoolInflationPayoutRate))
		Expect(updatedParams.MaxVotingPowerPerPool).To(Equal(math.LegacyMustNewDecFromStr("0.2")))
	})

	It("Update max voting power per pool with invalid value", func() {
		// ARRANGE
		payload := `{
			"max_voting_power_per_pool": "1.2"
		}`

		msg := &types.MsgUpdateParams{
			Authority: gov,
			Payload:   payload,
		}

		proposal, _ := govV1Types.NewMsgSubmitProposal(
			[]sdk.Msg{msg}, minDeposit, i.DUMMY[0], "", "title", "summary", false,
		)

		// ACT
		_, submitErr := s.RunTx(proposal)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		updatedParams := s.App().PoolKeeper.GetParams(s.Ctx())

		Expect(submitErr).To(HaveOccurred())

		Expect(updatedParams.ProtocolInflationShare).To(Equal(types.DefaultProtocolInflationShare))
		Expect(updatedParams.PoolInflationPayoutRate).To(Equal(types.DefaultPoolInflationPayoutRate))
		Expect(updatedParams.MaxVotingPowerPerPool).To(Equal(types.DefaultMaxVotingPowerPerPool))
	})
})
