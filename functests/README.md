# kubevirt-template-validator functional tests

*WARNING*: work in progress. We are still working to make the tests more automated and more robust.

## Preparation
1. install OKD >= 3.11
2. install kubevirt >= 0.14
4. install the common templates >= 0.4.1
3. install the webhook
5. ready!

## run the tests
```bash
./test-runner.sh
```

## TODO
- run the tests in a new namespace?
- remove (stale/bogus) reference to cirros to the functest manifests.
