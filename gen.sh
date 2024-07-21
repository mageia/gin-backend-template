PROTO_NAMES=(
    "hello"
)

for name in "${PROTO_NAMES[@]}"; do
  protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative,require_unimplemented_servers=false:. rpc/${name}.proto
  if [ $? -ne 0 ]; then
      echo "error processing ${name}.proto"
      exit $?
  fi
done
