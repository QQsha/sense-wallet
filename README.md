# sense-wallet

**Docker build:**
   ❯ docker build -t wallet .  
**Docker run:**
   ❯ docker run -p 8080:8080 wallet 
   
**Endpoints:**
  Make transaction(POST):
  http://localhost:8080/transaction
  {
	"user_id": "134256",
	"currency": "EUR",
	"amount": 0.5,
	"time_placed": "24-JAN-20 10:27:44",
	"type": "deposit"
}

  Get user balance(POST):
  http://localhost:8080/balance/134256
