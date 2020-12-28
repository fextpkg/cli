module github.com/Flacy/upip

go 1.15

require (
  github.com/go-ini/ini v1.62.0
  golang.org/x/net v0.0.0-20201224014010-6772e930b67b
  golang.org/x/sys v0.0.0-20201223074533-0d417f636930
)

replace github.com/Flacy/upip/base_cfg => ./base_cfg
replace github.com/Flacy/upip/utils => ./utils
replace github.com/Flacy/upip/color => ./color
replace github.com/Flacy/upip/whl => ./whl
replace github.com/Flacy/upip/cmd => ./cmd
replace github.com/Flacy/upip/help => ./help

