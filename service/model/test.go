package model

const TestCreateQuery = `
CREATE TABLE test (
id BIGINT(20) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
test_field VARCHAR(30) NOT NULL
)
`

func CreateTestTable() {

}
