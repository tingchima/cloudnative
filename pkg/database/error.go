package database

// // IsDuplicateErr return true if is duplicate error
// func IsDuplicateErr(err error) bool {
// 	if err == nil {
// 		return false
// 	}
// 	mysqlErr, ok := err.(*mysql.MySQLError)
// 	if ok {
// 		if mysqlErr.Number == 1062 {
// 			// solve the duplicate key error.
// 			return true
// 		}
// 	}
// 	return false
// }
