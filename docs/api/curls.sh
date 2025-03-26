#!/bin/bash

# transfer money
curl -X POST http://localhost:8088/api/sheba \
-H "Content-Type: application/json" \
-d '{"price": 50, "fromShebaNumber": "IR7740802513265426484548", "ToShebaNumber": "IR9470104877394934515563", "note": "توضیحات تراکنش"}'
