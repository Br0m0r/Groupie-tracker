Implement constant-time map lookups & optimized coordinate refresh

- repository/artist_repository.go  
  • Add `artistMap`, `locationMap`, `memberCountMap`, `creationYearMap`, `albumYearMap` and `minCreationYear`/`minAlbumYear` fields  
  • Populate all maps and minima in `LoadData()`  
  • Expose `GetArtistsByLocation`, `GetArtistsByMemberCount`, `GetArtistsByCreationYear`, `GetArtistsByAlbumYear`, `GetMinYears`

- repository/coordinates_repository.go  
  • Add `Has(location string) bool` and `CacheSize() int` helpers  
  • Preserve existing cache across swaps  

- store/store.go  
  • Add delegates `GetArtistsByMemberCount`, `GetArtistsByCreationYear`, `GetArtistsByAlbumYear`  
  • Revise `SwapData` to swap only `artistRepo`, keep `coordinatesRepo` intact, detect & fetch only new locations, and log “refresh started” + before→after artist/coord counts

- handlers/filter.go  
  • Add fast-path branches for pure-location, member-count, creation-year, and album-year filters using DataStore map lookups  
  • Retain fallback to full in-memory scan for mixed or free-text queries
