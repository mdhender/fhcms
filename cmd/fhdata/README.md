# fhdata

This command imports the `galaxy.dat`, `stars.dat`, `planets.dat` and `sp??.dat` files created by the 32-bit C program.

It creates two sets of output files.
The first set is a fairly faithful export of the original data files.
The names of these files are `galaxy.json`, `stars.json`, `planets.json`, and `sp??.json`.

The second set is actually a single file, `cluster.json`, and contains all the data merged into a single file.
This export also massages the data to make it easier to use by the `fhapp` command.

The merged data has the following structure:

1. Systems
   1. Planets
      1. Colonies
2. Species
   1. Named Planets
   2. Ships

Both colonies and ships have inventories stored with them.