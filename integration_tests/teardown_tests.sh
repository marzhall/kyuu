echo "Tearing down tests."
# wait for the server to stop handling any remaining read/write calls
sleep 1
killall kyuuD
# in the event kyuuD has failed, wait a moment for fusermount to see that
sleep 1
fusermount -u test_kyuus
rm -rf test_kyuus
rm -rf kyuu_storage
