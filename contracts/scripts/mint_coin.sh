root_path=$(git rev-parse --show-toplevel)
cd $root_path/contracts/testcoin
sui move build
sui client publish --gas-budget 1000000000 --skip-dependency-verification --json > mint_receipt.json
