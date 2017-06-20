rm -rf testing

mkdir testing
cd testing
codon init
cd ..

rsync -a compliance/ testing/

cd testing
make generate
# Patch clients with mock
rm -rf clients
cd ..
rsync -a compliance/clients/ testing/clients/

cd testing
go test

cd ..
rm -rf testing
