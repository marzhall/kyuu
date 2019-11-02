echo "Tearing down tests."
# wait for the server to stop handling any remaining read/write calls
sleep 1
fusermount -u test_kyuus
killall kyuuD
rm -rf test_kyuus
