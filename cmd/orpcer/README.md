# Orpcer

Orp[han]cer[ts] is looking for outdated certificates. Outdated means, you specify a threshold value and if the age (in days) of the cert is higher than your threshold is printing out the certificates.

## Purpose

The main purpose of orpcer is for running as cron. Each interval it checks for orphan certificates and print either "no orphan certs found" or the list of outdated certs. You can set up an alert which checks the logs that in your cron interval the string "no orphan certs found" appears at least once.
