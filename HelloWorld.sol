pragma solidity ^0.8.17;
contract HelloWorld {
    uint a;
    constructor() public {
      a = 0x55;
    }

    function sayHello() public pure returns(string memory) {
      return "helloworld";
    }
}