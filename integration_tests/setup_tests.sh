echo "Setting up test environment."
../bin/kyuuD &
sleep 1

if [ ! -d "test_kyuus" ]
then
    mkdir test_kyuus
else
    rm -rf test_kyuus
    mkdir test_kyuus
fi

9pfuse 127.0.0.1:5640 test_kyuus
