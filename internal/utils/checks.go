package utils

func IsStringSliceEmpty(input []string) bool {
  if len(input) > 0 {
    return false
  }

  return true
}

func IsStringEmpty(input string) bool {
  if len(input) > 0 {
    return false
  }

  return true
}

func IsIntEmpty(input int64) bool {
  if input != 0 {
    return false
  }

  return true
}
