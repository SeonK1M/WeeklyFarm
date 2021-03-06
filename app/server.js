// 모듈 포함
const express = require('express');
const app = express();
var bodyParser = require('body-parser');

// Constants
const PORT = 3000;
const HOST = "0.0.0.0";

 // 패브릭 연결설정
const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const ccpPath = path.resolve(__dirname, 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

app.use(express.static(path.join(__dirname, 'views')));
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

// Index page
app.get('/', (req, res) => {
    res.sendFile(__dirname + '/views/function.html');
})

// rest 라우팅
// url mate POST(type1 생성 type2 프로젝트 점수 추가)
app.post('/registerFarm', async(req, res) => {
    const fid = req.body.fid;
    const bid = req.body.bid;
    const address = req.body.address;
    const rentP = req.body.rentPrice;
    console.log("id: " + fid);

    // Wallet 가져오기
    const walletPath = path.join(process.cwd(), 'wallet');
    console.log("1111");
    const wallet = new FileSystemWallet(walletPath);
    console.log("2222");
    const userExists = await wallet.exists('user1');
    console.log("3333");
    if (!userExists) {
        console.log('No User1 in your wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    // Gateway에 연결하기
    const gateway = new Gateway();
    console.log("4444");
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false }});
    console.log("5555");
    // Channel에 연결하기
    const network = await gateway.getNetwork('mychannel');
    console.log("6666");
    // ChainCode 클래스 가져오기
    const contract = network.getContract('farm');
    console.log("7777");
    // submit transaction
    const result = await contract.submitTransaction("registerFarmland", fid, bid, address, rentP);

    // 체인코드 수행결과를 클라이언트에 알려주기
    console.log("well registered");
    res.status(200).send(`{result: "success", msg: "TX has been submitted"}`);
})

app.post('/newRecord', async(req, res) => {
    const cid = req.body.cid;
    const uid = req.body.uid;
    const fid = req.body.fid;
    const pname = req.body.pname;
    console.log("id is " + cid)

    // Wallet 가져오기
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);
    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    // Gateway에 연결하기
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false }});
    // Channel에 연결하기
    const network = await gateway.getNetwork('mychannel');
    // ChainCode 클래스 가져오기
    const contract = network.getContract('farm');
    // submit transaction
    const result = await contract.submitTransaction("newRecord", cid, uid, fid, pname);
    // 체인코스 수행결과를 클라이언트에 알려주기
    console.log("well recorded 'new record'");
    res.status(200).send(`{result: "success", msg: "TX has been submitted"}`);
})

app.post('/updateRecord', async(req, res) => {
    const cid = req.body.cid
    console.log("contract id is " + cid)

    // Wallet 가져오기
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);
    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    // Gateway에 연결하기
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false }});
    // Channel에 연결하기
    const network = await gateway.getNetwork('mychannel');
    // ChainCode 클래스 가져오기
    const contract = network.getContract('farm');
    // submit transaction
    const result = await contract.submitTransaction("updateRecord", cid);
    // 체인코스 수행결과를 클라이언트에 알려주기
    console.log("well recored 'update record'");
    res.status(200).send(`{result: "success", msg: "TX has been submitted"}`);
})

// // url mate GET(type1 조회 type2 이력조회)
app.get('/getRecord', async(req, res) => {
    const cid = req.query.cid;
    console.log(cid+"3333")
    // Wallet 가져오기
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);
    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    // Gateway에 연결하기
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false }});
    // Channel에 연결하기
    const network = await gateway.getNetwork('mychannel');
    // ChainCode 클래스 가져오기
    const contract = network.getContract('farm');
    // submit transactiond
    const result = await contract.evaluateTransaction("getRecordInfo", cid);
    // 체인코스 수행결과를 클라이언트에 알려주기
    console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
    res.status(200).send(result);
})

// server start
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);
