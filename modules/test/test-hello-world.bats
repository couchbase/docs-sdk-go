# #!./test/libs/bats/bin/bats

load 'test_helper'

@test "[hello-world] - startusing.go" {
    runExample $HELLO_WORLD_DIR startusing.go
    assert_success
    assert_output --partial "User: {Jade jade@test-email.com [Swimming Rowing]}"
    assert_output --partial "map[airline:map[callsign:MILE-AIR country:United States iata:Q5 icao:MLA id:10 name:40-Mile Air type:airline]]"
}
