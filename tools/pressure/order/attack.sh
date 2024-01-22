#!/usr/bin/env bash

vegeta attack -targets=targets.txt -output=results.bin -config=vegeta.conf

vegeta report -inputs=results.bin


