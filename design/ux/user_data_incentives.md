# User interaction with Obscuro 

In a typical transparent blockchain, there are service providers like "Infura" who act as the interface to the blockchain from the point of view of most users.

They can perform this service relatively cheap, because all data is visible to everyone, so they can employ efficient caching techniques. 
They also personally gain from understanding the user queries, and from analysing the exposed data.

As a result, the crypto users have become used to this as a free service, and take it for granted.

The data that is queried by users is:

1. Transaction receipts. Either queried after submitting a transaction to understand whether it was executed successfully and included in a block, 
   or to retrieve the submitted transactions.
2. RPC "eth_call" - which are used to query the status of the various smart contracts. For example, by far the most used: "ERC20.getBalance"   
3. RPC "get_balance" - to query the balance of the native ETH for an account.
4. Event query or subscription - used by dApps to update the UIs 


## Chains with privacy

It gets a bit more complicated when privacy is added to the mix, because the results of the user interactions are dependent on the requester.
Only logic run inside an approved TEE is able to return a result, since it is the only place where that data is visible.

This means that it is now much more difficult to create a cheap caching layer.
Also, the service provider will not gain any insights from this service since all traffic and processing is encrypted.


## Delegate to the nodes

The usual rescue in the crypto space is to employ incentives.

Obscuro is a decentralised network of nodes with different roles who already have their own incentives (Aggregators and Verifiers).

The ideal scenario is to have a large and diverse community of verifiers node to make sure that the network functions correctly.

Is is natural to assign to the verifier nodes the additional task of servicing user requests.

### Incentives

There has to be a reason for a node operator to provide this service to users. If there is no payoff, then most operators could just block the ports
and not accept incoming connections.

Here is a list and analyisis of various potential reasons for offering this service

#### 1. App developer

Someone developed a dApp and wants to support its users, since they are directly interested in the growth of the application. It could be a company, or a DAO.
E.g.: A game developer would not want to reply to requests made by a DEX user.

For this to be feasible, the enclave must be able to restrict requests to a whitelist of apps. Which is a technical feature that has to be implemented. 
Also the wallet extension (tooling) needs to know which node to use for which call.


#### 2. Fees

Node operators could charge fees for this service.

Given that everyone is now expecting this to be a free service, this is unlikely to be something that has a chance.


#### 3. Incentives payed by the protocol.

The network (or protocol) charges fees from user when submitting transactions. This is something that users expect to pay.

The Obscuro protocol is designed in such a way that it decouples the income from the costs by maintaing a buffer.

We can use this designed mechanism to pay for node usage as well along with the L1 gas fees and the general incentives to follow the protocol.


##### Measuring node usage

As a proxy for a node responding to user requests, we can use a model where a node is payed a percentage of the gas fees that originated from their node. 

A user that is connected to a node that doesn't respond to requests (like transaction receipts, or events), will leave that node and connect to a different node.

The "Wallet extension", our client-side software can hide this complexity from the user, and switch nodes without the user noticing.

This is not a perfect system, but it should be a good approximation. It turns a node operator into a service provider that makes an income based on the quality of service it provides. 

It also removes the need to incentivise verifiers to break the "verifier's dilemma"


##### Technical implementation

To implement the mechanism described above, the network needs to keep track of the original node that received the transaction via RPC.

Upon receiving a signed transaction from the wallet (MetaMask), the wallet extension will create a payload like this:

```
type RPCTxReceiver struct {
  SignedTransaction []byte   // The signed transaction generated by the wallet.
  RPCNode           Address  // The address of the node that accepted the request.
  R, S              *big.Int // Signature by the registered viewing key of the transaction signer.
}
```

The transaction is wrapped together with a node address, and signed with the viewing key. 

A node will only accept transactions via RPC if they're addressed to it, and thus can use this payload to claim a reward. 
This wrapper will then be gossiped to the other nodes accompanied by the signed viewing key.


##### Outstanding problems to solve

Problem 1: a malicious or malfunctioning wallet extension could submit the same transction wrapped for multiple nodes.
This is not a major problem because the incentives to do this are not large

Problem 2: how to include the transactions in the rollup. 
It is not ideal to include this extra wrapped information in the rollup and bulk it up. 

If the extra info is not included then it becomes harder to verify whether the fees were allcated correctly.
Nodes could also claim in bulk by constructing a ZKP that is included.


##### Payout

Assuming we solve the claiming problem, the question is how much should a node receive per transaction.

- A fixed amount per transaction
- A percentage of the fee paid by the transaction 

What amount would make running an L2 node that offers a useful service a reasonably profitable business?
Is that amount sustainable from the fees collected by the network? Is thre any help required from the foundation?


## Known problems 
- this mechanism creates the incentive for verifiers to DoS attack other verifiers, to get more business for themselves