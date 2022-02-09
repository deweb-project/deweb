# Storage

Data is stored in golang structs, that are stored in sqlite/gorm or other database format. For now binary data is also stored in the struct, this may have negative impact on the performance but that's not a problem now.

# File storage

Each conversation have a shared storage, which is based on git. All data is being commited and saved, then it is being pulled from the outside. That allow you to store history, and do some kind of merging.