ROOT_DIR="$(pwd)"

for goModFile in $(find . -name '*.mod'); do
  currentGoRootDirectory="$(dirname "$goModFile")"
  echo ""
  echo "### $currentGoRootDirectory ###"
  echo ""

  cd "$currentGoRootDirectory" || exit 1

  # synch packages and clear test cache
  go mod tidy
  go clean -testcache

  # run tests and create coverage report
  go test ./... -coverprofile=cover.out

  # check if all tests were successfull, otherwise return with error
  if [ $? -ne 0 ]; then
    rm -f cover.out
    exit 1
  fi

  # print only the total test coverage of the go module
  COVERAGE=$(
    go tool cover -func cover.out \
      | grep 'total' \
      | awk '{print $3}' \
  )

  rm -f cover.out

  COVERAGE_NUMBER=$(
    echo "${COVERAGE}" \
      | sed -e 's/\.[0-9]*//' -e 's/%//'
  )

  case ${COVERAGE_NUMBER} in
      [0-5][0-9]|[0-9]) COLOR="\033[31m" ;; # if coverage < 60%, then red color
      *) COLOR="\033[32m"                   # otherwise green
  esac

  echo "Total test coverage: ${COLOR}${COVERAGE}\033[0m"

  cd "$ROOT_DIR" || exit 1
done