package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"math/big"
	"strings"
	"time"
)

func newArguments(typeNames ...string) abi.Arguments {
	var args abi.Arguments
	for i, tn := range typeNames {
		abiType, err := abi.NewType(tn, tn, nil)
		if err != nil {
			panic(err)
		}
		args = append(args, abi.Argument{Name: fmt.Sprintf("%d", i), Type: abiType})
	}
	return args
}

func mustNewArguments(types ...string) (result abi.Arguments) {
	var err error
	for _, t := range types {
		var typ abi.Type
		items := strings.Split(t, " ")
		var name string
		if len(items) == 2 {
			name = items[1]
		} else {
			name = items[0]
		}
		typ, err = abi.NewType(items[0], items[0], nil)
		if err != nil {
			panic(err)
		}
		result = append(result, abi.Argument{Type: typ, Name: name})
	}
	return result
}

var verifyParliaBlockInput = mustNewArguments(
	"uint256 chainId",
	"bytes blockProof",
	"uint32 epochInterval",
)

var measureBlockInput = mustNewArguments(
	"bytes blockProof",
	"uint256 chainId",
)

// bytecode of VerifierGasMeasurer contract
func createBytecode() []byte {
	backend := backends.NewSimulatedBackend(core.GenesisAlloc{}, 8_000_000)
	bytecode, err := backend.CallContract(context.Background(), ethereum.CallMsg{
		Data: hexutil.MustDecode("0x60c060405234801561001057600080fd5b50600f60805260c860a05260805160a05161259f6100476000396000818161049f0152610e5f0152600061047e015261259f6000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c80632a36f0ad146100675780632b150b3c146100975780632dae36a0146100b7578063783d0962146100d757806386cd0354146100ea578063e4adf32f1461010a575b600080fd5b61007a610075366004611f7b565b61012b565b6040516001600160401b0390911681526020015b60405180910390f35b6100aa6100a5366004611fc6565b61014f565b60405161008e9190612007565b6100ca6100c53660046120ec565b610387565b60405161008e9190612180565b6100ca6100e5366004611f7b565b6103f4565b6100fd6100f8366004611f7b565b610464565b60405161008e91906121e0565b61011d6101183660046122bf565b610477565b60405161008e929190612338565b60005a905061013b848484610842565b505a6101479082612379565b949350505050565b604080516101a081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e0810182905261010081018290526101208101829052610140810182905261016081018290526101808101829052906101c28484610fb6565b90506101cd81610fcc565b82526101d881610fdf565b90506101e381610fcc565b60208301526101f181610fdf565b90506101fc81610ff4565b6001600160a01b0316604083015261021381610fdf565b905061021e81610fcc565b606083015261022c81610fdf565b905061023781610fcc565b608083015261024581610fdf565b905061025081610fcc565b60a083015261025e81610fdf565b905061026981610fdf565b905061027481610fdf565b90506102888161028383611001565b6110a4565b6001600160401b031660c083015261029f81610fdf565b90506102ae8161028383611001565b6001600160401b031660e08301526102c581610fdf565b90506102d48161028383611001565b6001600160401b03166101008301526102ec81610fdf565b90506102fb8161028383611001565b6001600160401b031661012083015261031381610fdf565b905061031e81610fdf565b905061032981610fcc565b61014083015261033881610fdf565b90506103478161028383611001565b6001600160401b031661016083015261035f81610fdf565b905083836040516103719291906123a1565b6040519081900390206101808301525092915050565b60606000610396868686610842565b805190915083146103e75760405162461bcd60e51b81526020600482015260166024820152756e6f74206120636865636b706f696e7420626c6f636b60501b60448201526064015b60405180910390fd5b6060015195945050505050565b60606000610403858585610842565b905080602001516001600160401b03166000146104585760405162461bcd60e51b81526020600482015260136024820152726e6f7420612067656e6573697320626c6f636b60681b60448201526064016103de565b60600151949350505050565b61046c611efc565b610147848484610842565b60606000807f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000063ffffffff82168910156105125760405162461bcd60e51b815260206004820152601b60248201527a0dac2d6ca40e6eae4ca40e0e4dedecce640c2e4ca40cadcdeeaced602b1b60448201526064016103de565b6000896001600160401b0381111561052c5761052c6123b1565b604051908082528060200260200182016040528015610555578160200160208202803683370190505b5090506000805b8463ffffffff168110156107de5760006105998e8e84818110610581576105816123c7565b905060200281019061059391906123dd565b8e610842565b60408101519091506000805b8c81101561060757826001600160a01b03168e8e838181106105c9576105c96123c7565b90506020020160208101906105de9190612423565b6001600160a01b0316036105f55760019150610607565b806105ff8161244c565b9150506105a5565b50806106465760405162461bcd60e51b815260206004820152600e60248201526d3ab735b737bbb71039b4b3b732b960911b60448201526064016103de565b6000805b866001600160401b03168110156106a857836001600160a01b0316888281518110610677576106776123c7565b60200260200101516001600160a01b03160361069657600191506106a8565b806106a08161244c565b91505061064a565b50806106f0578287876001600160401b0316815181106106ca576106ca6123c7565b6001600160a01b0390921660209283029190910190910152856106ec81612465565b9650505b8460000361077d578763ffffffff16846020015161070e91906124a1565b6001600160401b0316156107525760405162461bcd60e51b815260206004820152600b60248201526a65706f636820626c6f636b60a81b60448201526064016103de565b8763ffffffff16846020015161076891906124c7565b9a5083606001519b50836000015199506107c7565b898460800151146107c25760405162461bcd60e51b815260206004820152600f60248201526e0c4c2c840e0c2e4cadce840d0c2e6d608b1b60448201526064016103de565b835199505b5050505080806107d69061244c565b91505061055c565b508363ffffffff16816001600160401b031610156108335760405162461bcd60e51b81526020600482015260126024820152711c5d5bdc9d5b481b9bdd081c995858da195960721b60448201526064016103de565b50505050509550959350505050565b61084a611efc565b61ffff83111561085957600080fd5b83600061086582611001565b9050610870826110b7565b915061087b82610fcc565b608084015261088c6101bd836124ed565b915061089782610fdf565b91506108a68261028384611001565b6001600160401b031660208401526108c76108c2808085610fdf565b610fdf565b9150816108d381610fdf565b92508260006108e1836110c2565b90506000806041836108f38787612505565b6108fd9190612505565b6109079190612505565b9050603881101561091b5760019150610932565b610924816110cd565b61092f9060016124ed565b91505b50600061093e89611160565b905060006041838589855161095391906124ed565b61095d9190612505565b61096791906124ed565b6109719190612505565b6001600160401b03811115610988576109886123b1565b6040519080825280601f01601f1916602001820160405280156109b2576020820181803683370190505b50905060f960f81b816000815181106109cd576109cd6123c7565b60200101906001600160f81b031916908160001a9053506000600382516109f49190612505565b9050600881901c60f81b82600181518110610a1157610a116123c7565b60200101906001600160f81b031916908160001a905350600081901c60f81b82600281518110610a4357610a436123c7565b60200101906001600160f81b031916908160001a9053505060005b8251811015610acd57828181518110610a7957610a796123c7565b01602001516001600160f81b03191682610a948360036124ed565b81518110610aa457610aa46123c7565b60200101906001600160f81b031916908160001a90535080610ac58161244c565b915050610a5e565b506023825101810160038d01808803808284379190910184019050868501604019888803879003018082843791909101905085602a8082843750505060008060418689890303039150602a85830101835103905084600403610c2257605d60f91b83610b3a8360006124ed565b81518110610b4a57610b4a6123c7565b60200101906001600160f81b031916908160001a9053506001600160f81b031960e883901b1683610b7c8360016124ed565b81518110610b8c57610b8c6123c7565b60200101906001600160f81b031916908160001a9053506001600160f81b031960f083901b1683610bbe8360026124ed565b81518110610bce57610bce6123c7565b60200101906001600160f81b031916908160001a90535060f882901b83610bf68360036124ed565b81518110610c0657610c066123c7565b60200101906001600160f81b031916908160001a905350610d5b565b84600303610cb55760b960f81b83610c3b8360006124ed565b81518110610c4b57610c4b6123c7565b60200101906001600160f81b031916908160001a9053506001600160f81b031960f083901b1683610c7d8360016124ed565b81518110610c8d57610c8d6123c7565b60200101906001600160f81b031916908160001a90535060f882901b83610bf68360026124ed565b84600203610d1057601760fb1b83610cce8360006124ed565b81518110610cde57610cde6123c7565b60200101906001600160f81b031916908160001a9053506001600160f81b031960f083901b1683610bf68360016124ed565b6038821015610d5b57610d248260806124ed565b60f81b83610d338360006124ed565b81518110610d4357610d436123c7565b60200101906001600160f81b031916908160001a9053505b5050604080516041808252608082019092526000916020820181803683370190505090506041808703602083013760208a01516001600160401b031615610e5d578051600160f81b9082906040908110610db757610db76123c7565b01602001516001600160f81b03191603610dff57601c60f81b81604081518110610de357610de36123c7565b60200101906001600160f81b031916908160001a905350610e2f565b601b60f81b81604081518110610e1757610e176123c7565b60200101906001600160f81b031916908160001a9053505b60a08a0182905260c08a0181905281516020830120610e4e9082611aec565b6001600160a01b031660408b01525b7f000000000000000000000000000000000000000000000000000000000000000063ffffffff168a60200151610e9391906124a1565b6001600160401b0316600003610f8b57600060146020604188610eb68c8c612505565b610ec091906124ed565b610eca9190612505565b610ed49190612505565b610ede919061251c565b90506000816001600160401b03811115610efa57610efa6123b1565b604051908082528060200260200182016040528015610f23578160200160208202803683370190505b50905060005b82811015610f835760006020601483028a8d010101359050606081901c838381518110610f5857610f586123c7565b6001600160a01b03909216602092830291909101909101525080610f7b8161244c565b915050610f29565b5060608c0152505b8c8c604051610f9b9291906123a1565b6040519081900390208a525050505050505050509392505050565b600082610fc281611b10565b61014790826124ed565b6000610fd9826021611b8a565b92915050565b6000610fea82611001565b610fd990836124ed565b6000610fd9826015611b8a565b6000808235811a608081101561101a576001915061109d565b60b88110156110405761102e608082612505565b6110399060016124ed565b915061109d565b60c081101561106b576001939093019283356008602083900360b701021c810160b51901915061109d565b60f881101561107f5761102e60c082612505565b6001939093019283356008602083900360f701021c810160f5190191505b5092915050565b60006110b08383611b8a565b9392505050565b6000610fea82611b10565b6000610fd982611b10565b60006101008210156110e157506001919050565b620100008210156110f457506002919050565b630100000082101561110857506003919050565b600160201b82101561111c57506004919050565b600160281b82101561113057506005919050565b600160301b82101561114457506006919050565b600160381b82101561115857506007919050565b506008919050565b6060816000036111bd5760408051600180825281830190925290602082018180368337019050509050608060f81b816000815181106111a1576111a16123c7565b60200101906001600160f81b031916908160001a905350919050565b607f82116111fb57604080516001808252818301909252906020820181803683370190505090508160f81b816000815181106111a1576111a16123c7565b61010082101561126a5760408051600280825281830190925290602082018180368337019050509050608160f81b8160008151811061123c5761123c6123c7565b60200101906001600160f81b031916908160001a9053508160f81b816001815181106111a1576111a16123c7565b6201000082101561130c5760408051600380825281830190925290602082018180368337019050509050608260f81b816000815181106112ac576112ac6123c7565b60200101906001600160f81b031916908160001a905350600882901c60f81b816001815181106112de576112de6123c7565b60200101906001600160f81b031916908160001a9053508160f81b816002815181106111a1576111a16123c7565b63010000008210156113e15760408051600480825281830190925290602082018180368337019050509050608360f81b8160008151811061134f5761134f6123c7565b60200101906001600160f81b031916908160001a905350601082901c60f81b81600181518110611381576113816123c7565b60200101906001600160f81b031916908160001a905350600882901c60f81b816002815181106113b3576113b36123c7565b60200101906001600160f81b031916908160001a9053508160f81b816003815181106111a1576111a16123c7565b600160201b8210156114e85760408051600580825281830190925290602082018180368337019050509050608460f81b81600081518110611424576114246123c7565b60200101906001600160f81b031916908160001a905350601882901c60f81b81600181518110611456576114566123c7565b60200101906001600160f81b031916908160001a905350601082901c60f81b81600281518110611488576114886123c7565b60200101906001600160f81b031916908160001a905350600882901c60f81b816003815181106114ba576114ba6123c7565b60200101906001600160f81b031916908160001a9053508160f81b816004815181106111a1576111a16123c7565b600160281b8210156116215760408051600680825281830190925290602082018180368337019050509050608560f81b8160008151811061152b5761152b6123c7565b60200101906001600160f81b031916908160001a905350602082901c60f81b8160018151811061155d5761155d6123c7565b60200101906001600160f81b031916908160001a905350601882901c60f81b8160028151811061158f5761158f6123c7565b60200101906001600160f81b031916908160001a905350601082901c60f81b816003815181106115c1576115c16123c7565b60200101906001600160f81b031916908160001a905350600882901c60f81b816004815181106115f3576115f36123c7565b60200101906001600160f81b031916908160001a9053508160f81b816005815181106111a1576111a16123c7565b600160301b82101561178c5760408051600780825281830190925290602082018180368337019050509050608660f81b81600081518110611664576116646123c7565b60200101906001600160f81b031916908160001a905350602882901c60f81b81600181518110611696576116966123c7565b60200101906001600160f81b031916908160001a905350602082901c60f81b816002815181106116c8576116c86123c7565b60200101906001600160f81b031916908160001a905350601882901c60f81b816003815181106116fa576116fa6123c7565b60200101906001600160f81b031916908160001a905350601082901c60f81b8160048151811061172c5761172c6123c7565b60200101906001600160f81b031916908160001a905350600882901c60f81b8160058151811061175e5761175e6123c7565b60200101906001600160f81b031916908160001a9053508160f81b816006815181106111a1576111a16123c7565b600160381b8210156119295760408051600880825281830190925290602082018180368337019050509050608760f81b816000815181106117cf576117cf6123c7565b60200101906001600160f81b031916908160001a905350603082901c60f81b81600181518110611801576118016123c7565b60200101906001600160f81b031916908160001a905350602882901c60f81b81600281518110611833576118336123c7565b60200101906001600160f81b031916908160001a905350602082901c60f81b81600381518110611865576118656123c7565b60200101906001600160f81b031916908160001a905350601882901c60f81b81600481518110611897576118976123c7565b60200101906001600160f81b031916908160001a905350601082901c60f81b816005815181106118c9576118c96123c7565b60200101906001600160f81b031916908160001a905350600882901c60f81b816006815181106118fb576118fb6123c7565b60200101906001600160f81b031916908160001a9053508160f81b816007815181106111a1576111a16123c7565b60408051600980825281830190925290602082018180368337019050509050608860f81b81600081518110611960576119606123c7565b60200101906001600160f81b031916908160001a905350603882901c60f81b81600181518110611992576119926123c7565b60200101906001600160f81b031916908160001a905350603082901c60f81b816002815181106119c4576119c46123c7565b60200101906001600160f81b031916908160001a905350602882901c60f81b816003815181106119f6576119f66123c7565b60200101906001600160f81b031916908160001a905350602082901c60f81b81600481518110611a2857611a286123c7565b60200101906001600160f81b031916908160001a905350601882901c60f81b81600581518110611a5a57611a5a6123c7565b60200101906001600160f81b031916908160001a905350601082901c60f81b81600681518110611a8c57611a8c6123c7565b60200101906001600160f81b031916908160001a905350600882901c60f81b81600781518110611abe57611abe6123c7565b60200101906001600160f81b031916908160001a9053508160f81b816008815181106111a1576111a16123c7565b6000806000611afb8585611bc8565b91509150611b0881611c36565b509392505050565b60008135811a6080811015611b285750600092915050565b60b8811080611b43575060c08110801590611b43575060f881105b15611b515750600192915050565b60c0811015611b7e57611b66600160b8612530565b611b739060ff1682612505565b6110b09060016124ed565b611b66600160f8612530565b60008082118015611b9c575060218211155b611ba557600080fd5b6000611bb084611b10565b93840135939092036020036008029290921c92915050565b6000808251604103611bfe5760208301516040840151606085015160001a611bf287828585611dea565b94509450505050611c2f565b8251604003611c275760208301516040840151611c1c868383611ecd565b935093505050611c2f565b506000905060025b9250929050565b6000816004811115611c4a57611c4a612553565b03611c525750565b6001816004811115611c6657611c66612553565b03611cae5760405162461bcd60e51b815260206004820152601860248201527745434453413a20696e76616c6964207369676e617475726560401b60448201526064016103de565b6002816004811115611cc257611cc2612553565b03611d0f5760405162461bcd60e51b815260206004820152601f60248201527f45434453413a20696e76616c6964207369676e6174757265206c656e6774680060448201526064016103de565b6003816004811115611d2357611d23612553565b03611d7b5760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c604482015261756560f01b60648201526084016103de565b6004816004811115611d8f57611d8f612553565b03611de75760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202776272076616c604482015261756560f01b60648201526084016103de565b50565b6000806fa2a8918ca85bafe22016d0b997e4df60600160ff1b03831115611e175750600090506003611ec4565b8460ff16601b14158015611e2f57508460ff16601c14155b15611e405750600090506004611ec4565b6040805160008082526020820180845289905260ff881692820192909252606081018690526080810185905260019060a0016020604051602081039080840390855afa158015611e94573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b038116611ebd57600060019250925050611ec4565b9150600090505b94509492505050565b6000806001600160ff1b03831660ff84901c601b01611eee87828885611dea565b935093505050935093915050565b6040805160e0810182526000808252602082018190529181018290526060808201819052608082019290925260a0810182905260c081019190915290565b60008083601f840112611f4c57600080fd5b5081356001600160401b03811115611f6357600080fd5b602083019150836020828501011115611c2f57600080fd5b600080600060408486031215611f9057600080fd5b83356001600160401b03811115611fa657600080fd5b611fb286828701611f3a565b909790965060209590950135949350505050565b60008060208385031215611fd957600080fd5b82356001600160401b03811115611fef57600080fd5b611ffb85828601611f3a565b90969095509350505050565b60006101a0820190508251825260208301516020830152604083015161203860408401826001600160a01b03169052565b50606083015160608301526080830151608083015260a083015160a083015260c083015161207160c08401826001600160401b03169052565b5060e083015161208c60e08401826001600160401b03169052565b50610100838101516001600160401b038116848301525050610120838101516001600160401b0381168483015250506101408381015190830152610160808401516001600160401b03811682850152505061018092830151919092015290565b6000806000806060858703121561210257600080fd5b84356001600160401b0381111561211857600080fd5b61212487828801611f3a565b90989097506020870135966040013595509350505050565b600081518084526020808501945080840160005b838110156121755781516001600160a01b031687529582019590820190600101612150565b509495945050505050565b6020815260006110b0602083018461213c565b6000815180845260005b818110156121b95760208185018101518683018201520161219d565b818111156121cb576000602083870101525b50601f01601f19169290920160200192915050565b60208152815160208201526001600160401b03602083015116604082015260018060a01b0360408301511660608201526000606083015160e0608084015261222c61010084018261213c565b9050608084015160a084015260a0840151601f19808584030160c08601526122548383612193565b925060c08601519150808584030160e0860152506122728282612193565b95945050505050565b60008083601f84011261228d57600080fd5b5081356001600160401b038111156122a457600080fd5b6020830191508360208260051b8501011115611c2f57600080fd5b6000806000806000606086880312156122d757600080fd5b85356001600160401b03808211156122ee57600080fd5b6122fa89838a0161227b565b909750955060208801359450604088013591508082111561231a57600080fd5b506123278882890161227b565b969995985093965092949392505050565b60408152600061234b604083018561213c565b90506001600160401b03831660208301529392505050565b634e487b7160e01b600052601160045260246000fd5b60006001600160401b038381169083168181101561239957612399612363565b039392505050565b8183823760009101908152919050565b634e487b7160e01b600052604160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b6000808335601e198436030181126123f457600080fd5b8301803591506001600160401b0382111561240e57600080fd5b602001915036819003821315611c2f57600080fd5b60006020828403121561243557600080fd5b81356001600160a01b03811681146110b057600080fd5b60006001820161245e5761245e612363565b5060010190565b60006001600160401b0380831681810361248157612481612363565b6001019392505050565b634e487b7160e01b600052601260045260246000fd5b60006001600160401b03808416806124bb576124bb61248b565b92169190910692915050565b60006001600160401b03808416806124e1576124e161248b565b92169190910492915050565b6000821982111561250057612500612363565b500190565b60008282101561251757612517612363565b500390565b60008261252b5761252b61248b565b500490565b600060ff821660ff84168082101561254a5761254a612363565b90039392505050565b634e487b7160e01b600052602160045260246000fdfea264697066735822122086a7c9dc7f5ed782523040256f09b30a09ff0a30b85c14b13e3b5fc40f5d575f64736f6c634300080e0033"),
	}, nil)
	if err != nil {
		panic(err)
	}
	return bytecode
}

var parliaBlockVerificationBytecode = createBytecode()

func main() {
	to := common.HexToAddress("0x0000000000000000000000000000000000000001")
	backend := backends.NewSimulatedBackend(core.GenesisAlloc{
		to: {
			Code:    parliaBlockVerificationBytecode,
			Balance: big.NewInt(0),
		},
	}, 8_000_000)
	blockProof := hexutil.MustDecode("0xf903fca00000000000000000000000000000000000000000000000000000000000000000a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d4934794fffffffffffffffffffffffffffffffffffffffea0919fcc7ad870b53db0aa76eb588da06bacb6d230195100699fc928511003b422a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001808402625a0080845e9da7ceb9020500000000000000000000000000000000000000000000000000000000000000002a7cdd959bfe8d9487b2a43b33565295a698f7e26488aa4d1955ee33403f8ccb1d4de5fb97c7ade29ef9f4360c606c7ab4db26b016007d3ad0ab86a0ee01c3b1283aa067c58eab4709f85e99d46de5fe685b1ded8013785d6623cc18d214320b6bb6475978f3adfc719c99674c072166708589033e2d9afec2be4ec20253b8642161bc3f444f53679c1f3d472f7be8361c80a4c1e7e9aaf001d0877f1cfde218ce2fd7544e0b2cc94692d4a704debef7bcb61328b8f7166496996a7da21cf1f1b04d9b3e26a3d0772d4c407bbe49438ed859fe965b140dcf1aab71a96bbad7cf34b5fa511d8e963dbba288b1960e75d64430b3230294d12c6ab2aac5c2cd68e80b16b581ea0a6e3c511bbd10f4519ece37dc24887e11b55d7ae2f5b9e386cd1b50a4550696d957cb4900f03a82012708dafc9e1b880fd083b32182b869be8e0922b81f8e175ffde54d797fe11eb03f9e3bf75f1d68bf0b8b6fb4e317a0f9d6f03eaf8ce6675bc60d8c4d90829ce8f72d0163c1d5cf348a862d55063035e7a025f4da968de7e4d7e4004197917f4070f1d6caa02bbebaebb5d7e581e4b66559e635f805ff0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000880000000000000000")
	const rounds = 100_000
	// solidity call
	input, err := measureBlockInput.Pack(blockProof, big.NewInt(65))
	if err != nil {
		panic(err)
	}
	solidityTime := time.Now().UnixNano()
	var gasUsed uint64
	for i := 0; i < rounds; i++ {
		gasUsed, err = backend.EstimateGas(context.Background(), ethereum.CallMsg{
			To:   &to,
			Data: append(hexutil.MustDecode("0x2a36f0ad"), input...),
		})
		if err != nil {
			panic(err)
		}
	}
	solidityTime = time.Now().UnixNano() - solidityTime
	fmt.Printf("elapsed time (solidity): %d ns\n", solidityTime)
	fmt.Printf("gas used (solidity): %d\n", gasUsed)
	// native call
	contract := vm.NewVerifyParliaBlock()
	input, err = verifyParliaBlockInput.Pack(big.NewInt(56), blockProof, uint32(200))
	if err != nil {
		panic(err)
	}
	nativeTime := time.Now().UnixNano()
	for i := 0; i < rounds; i++ {
		_, err = contract.Run(input)
		if err != nil {
			panic(err)
		}
	}
	nativeTime = time.Now().UnixNano() - nativeTime
	fmt.Printf("elapsed time (native): %d ns\n", nativeTime)
	fmt.Printf("optimal gas (native): ~%f\n", float64(gasUsed)*float64(nativeTime)/float64(solidityTime))
}