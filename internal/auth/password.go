package auth

import "golang.org/x/crypto/bcrypt"

// PasswordManager handles password hashing and verification
type PasswordManager struct{}

func NewPasswordManager() *PasswordManager {
	return &PasswordManager{}
}




// HashPassword hashes a plain text password using bcrypt
// working :- 
//1. Convert password → bytes
//bcrypt generates random salt
//Combine Password + Salt
//bcrypt hashing (multiple rounds) --> 2^Cost rounds of hashing
//returns hashed password string
// generated hash contain its salt + cost  [ final hash --> algorithm + cost + salt + hash ]
func (m *PasswordManager) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}



// VerifyPassword checks if a plain text password matches a hash
// working :-
//Extract salt and cost from stored hash
//Re-hash entered password using same salt + cost
//Compare new hash with stored hash
func (m *PasswordManager) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// advantage
//Configurable cost factor → makes brute-force attacks much harder by increasing computation time. 
// (high cost -> more round of hasing (slow to generate) -> more secure)
//generate automatic random salt to which make it more secure