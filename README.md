# Logmatch

application used to test <https://github.com/sjiekak/logen>

Using it with json logging applications

```
cat someapp.log | while IFS= read -r line; do jq -r .message ; done | go run github.com/sjiekak/logmatch
```

eg on an elasticsearch instance the relevant output look like

```log
338 matches with event "PrimaryreplicaResyncCompletedWithOperations" for line "primary-replica resync completed with 0 operations"
78 matches with event "LoadedModule" for line "loaded module [rest-root]"
42 matches with event "MemoryUsageDownAfter" for line "memory usage down after [0], before [123456789], after [12345]"
17 matches with event "Overhead" for line "[gc][906] overhead, spent [316ms] collecting in the last [1s]"
10 matches with event "GcDidBringMemoryUsageDown" for line "GC did bring memory usage down, before [123456789], after [12345], allocations [96], duration [296]"
10 matches with event "AttemptingToTriggerGgcDueToHighHeapUsage" for line "attempting to trigger G1GC due to high heap usage [123456789]"
4 matches with event "HandlingRequest" for line "handling request [InboundMessage{Header{333}{4444}{11111111}{true}{false}{false}{false}{indices:data/read/search[free_context/scroll]}}] took [5412ms] which is above the warn threshold of [5000ms]"
3 matches with event "FinishedWithResponseBulkByScrollResponse" for line "249048415 finished with response BulkByScrollResponse[took=1ms,timed_out=false,sliceId=null,updated=0,created=0,deleted=0,batches=0,versionConflicts=0,noops=0,retries=0,throttledUntil=0s,bulk_failures=[],search_failures=[]]"
```

## Ideas

- [Better match using levenshtein distance](https://github.com/agnivade/levenshtein)
