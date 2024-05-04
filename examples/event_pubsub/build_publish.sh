root_path=$(git rev-parse --show-toplevel)
cd $root_path/examples/event_pubsub/eventpub
sui move build
sui client publish --gas-budget 1000000000 --skip-dependency-verification --json > publish_receipt.json
