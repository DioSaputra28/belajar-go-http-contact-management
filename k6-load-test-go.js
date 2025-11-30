import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');

// Configuration
const BASE_URL = 'http://go-compare.diosaputra.site';

// Load test configuration with realistic traffic pattern
export const options = {
  stages: [
    // Warm up - gradually increase load
    { duration: '2m', target: 50 },    // Ramp up to 50 users over 2 minutes
    { duration: '3m', target: 100 },   // Ramp up to 100 users over 3 minutes
    
    // Peak hours - simulate busy periods
    { duration: '5m', target: 500 },   // Ramp up to 500 users (peak traffic)
    { duration: '5m', target: 500 },   // Stay at 500 users for 5 minutes
    
    // Spike test - sudden traffic increase
    { duration: '1m', target: 1000 },  // Spike to 1000 users
    { duration: '3m', target: 1000 },  // Maintain spike for 3 minutes
    
    // Cool down - traffic decreases
    { duration: '2m', target: 500 },   // Drop to 500 users
    { duration: '3m', target: 200 },   // Drop to 200 users
    { duration: '2m', target: 100 },   // Drop to 100 users
    
    // Low traffic period
    { duration: '3m', target: 50 },    // Drop to 50 users
    { duration: '2m', target: 0 },     // Ramp down to 0 users
  ],
  
  thresholds: {
    http_req_duration: ['p(95)<2000'], // 95% of requests should be below 2s
    http_req_failed: ['rate<0.1'],     // Error rate should be less than 10%
    errors: ['rate<0.1'],              // Custom error rate should be less than 10%
  },
};

// Generate random data
function randomString(length) {
  const chars = 'abcdefghijklmnopqrstuvwxyz';
  let result = '';
  for (let i = 0; i < length; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return result;
}

function randomEmail() {
  return `user_${randomString(8)}_${Date.now()}@test.com`;
}

function randomPhone() {
  return `08${Math.floor(Math.random() * 1000000000).toString().padStart(10, '0')}`;
}

// Main test scenario
export default function () {
  const email = randomEmail();
  const password = 'password123';
  let token = '';
  let userId = '';
  let contactId = '';
  let addressId = '';

  // 1. Register User
  {
    const payload = JSON.stringify({
      name: `Test User ${randomString(5)}`,
      email: email,
      password: password,
    });

    const params = {
      headers: { 'Content-Type': 'application/json' },
    };

    const res = http.post(`${BASE_URL}/user`, payload, params);
    
    const success = check(res, {
      'register status is 201': (r) => r.status === 201,
      'register has message': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.message === 'User created successfully';
        } catch (e) {
          return false;
        }
      },
    });

    errorRate.add(!success);

    // Debug: log response if failed
    if (!success) {
      console.log(`Register failed: ${res.status} - ${res.body}`);
    }
  }

  sleep(2); // Increased delay to ensure DB commit

  // 2. Login (test authentication)
  {
    const payload = JSON.stringify({
      email: email,
      password: password,
    });

    const params = {
      headers: { 'Content-Type': 'application/json' },
    };

    const res = http.post(`${BASE_URL}/login`, payload, params);
    
    const success = check(res, {
      'login status is 200': (r) => r.status === 200,
      'login returns token': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.user && body.user.token;
        } catch (e) {
          return false;
        }
      },
    });

    errorRate.add(!success);

    // Debug: log response if failed
    if (!success) {
      console.log(`Login failed for ${email}: ${res.status} - ${res.body}`);
    }

    if (success && res.status === 200) {
      const body = JSON.parse(res.body);
      token = body.user.token;
      userId = body.user.user_id;
    }
  }

  sleep(1);

  // 3. Get All Users
  if (token) {
    const params = {
      headers: {
        'Authorization': token,
      },
    };

    const res = http.get(`${BASE_URL}/user`, params);
    
    const success = check(res, {
      'get users status is 200': (r) => r.status === 200,
      'get users returns data': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.data && Array.isArray(body.data);
        } catch (e) {
          return false;
        }
      },
    });

    errorRate.add(!success);
  }

  sleep(1);

  // 4. Get User by ID
  if (token && userId) {
    const params = {
      headers: {
        'Authorization': token,
      },
    };

    const res = http.get(`${BASE_URL}/user/${userId}`, params);
    
    const success = check(res, {
      'get user by id status is 200': (r) => r.status === 200,
      'get user by id returns data': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.data && body.data.user_id;
        } catch (e) {
          return false;
        }
      },
    });

    errorRate.add(!success);
  }

  sleep(1);

  // 5. Create Contact
  if (token) {
    const payload = JSON.stringify({
      first_name: `First${randomString(5)}`,
      last_name: `Last${randomString(5)}`,
      email: randomEmail(),
      phone: randomPhone(),
    });

    const params = {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': token,
      },
    };

    const res = http.post(`${BASE_URL}/contact`, payload, params);
    
    const success = check(res, {
      'create contact status is 201': (r) => r.status === 201,
      'create contact returns data': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.message === 'Contact created successfully';
        } catch (e) {
          return false;
        }
      },
    });

    errorRate.add(!success);

    // Note: Go API doesn't return contact_id in response, need to fetch it
  }

  sleep(1);

  // 6. Get All Contacts
  if (token) {
    const params = {
      headers: {
        'Authorization': token,
      },
    };

    const res = http.get(`${BASE_URL}/contact`, params);
    
    const success = check(res, {
      'get contacts status is 200': (r) => r.status === 200,
      'get contacts returns array': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.data && Array.isArray(body.data);
        } catch (e) {
          return false;
        }
      },
    });

    errorRate.add(!success);

    // Get the first contact ID for further operations
    if (success && res.status === 200) {
      const body = JSON.parse(res.body);
      if (body.data && body.data.length > 0) {
        contactId = body.data[0].contact_id;
      }
    }
  }

  sleep(1);

  // 7. Get Contact by ID
  if (token && contactId) {
    const params = {
      headers: {
        'Authorization': token,
      },
    };

    const res = http.get(`${BASE_URL}/contact/${contactId}`, params);
    
    const success = check(res, {
      'get contact by id status is 200': (r) => r.status === 200,
      'get contact by id returns data': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.data && body.data.contact_id;
        } catch (e) {
          return false;
        }
      },
    });

    errorRate.add(!success);
  }

  sleep(1);

  // 8. Update Contact
  if (token && contactId) {
    const payload = JSON.stringify({
      first_name: `Updated${randomString(5)}`,
      last_name: `Last${randomString(5)}`,
      email: randomEmail(),
      phone: randomPhone(),
    });

    const params = {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': token,
      },
    };

    const res = http.put(`${BASE_URL}/contact/${contactId}`, payload, params);
    
    const success = check(res, {
      'update contact status is 200': (r) => r.status === 200,
      'update contact returns updated data': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.message === 'Contact updated successfully';
        } catch (e) {
          return false;
        }
      },
    });

    errorRate.add(!success);
  }

  sleep(1);

  // 9. Create Address
  if (token && contactId) {
    const cities = ['Jakarta', 'Bandung', 'Surabaya', 'Medan', 'Semarang'];
    const provinces = ['DKI Jakarta', 'Jawa Barat', 'Jawa Timur', 'Sumatera Utara', 'Jawa Tengah'];
    const idx = Math.floor(Math.random() * cities.length);

    const payload = JSON.stringify({
      street: `Jl. ${randomString(10)} No. ${Math.floor(Math.random() * 100)}`,
      city: cities[idx],
      province: provinces[idx],
      country: 'Indonesia',
      postal_code: `${Math.floor(Math.random() * 90000) + 10000}`,
      contact_id: contactId,
    });

    const params = {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': token,
      },
    };

    const res = http.post(`${BASE_URL}/address/`, payload, params);
    
    const success = check(res, {
      'create address status is 201': (r) => r.status === 201,
      'create address returns data': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.message === 'Address created successfully';
        } catch (e) {
          return false;
        }
      },
    });

    errorRate.add(!success);
  }

  sleep(1);

  // 10. Get All Addresses
  if (token && contactId) {
    const params = {
      headers: {
        'Authorization': token,
      },
    };

    const res = http.get(`${BASE_URL}/address/${contactId}`, params);
    
    const success = check(res, {
      'get addresses status is 200': (r) => r.status === 200,
      'get addresses returns array': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.data && Array.isArray(body.data);
        } catch (e) {
          return false;
        }
      },
    });

    errorRate.add(!success);

    // Get the first address ID for further operations
    if (success && res.status === 200) {
      const body = JSON.parse(res.body);
      if (body.data && body.data.length > 0) {
        addressId = body.data[0].address_id;
      }
    }
  }

  sleep(1);

  // 11. Get Address by ID
  if (token && contactId && addressId) {
    const params = {
      headers: {
        'Authorization': token,
      },
    };

    const res = http.get(`${BASE_URL}/address/${contactId}/${addressId}`, params);
    
    const success = check(res, {
      'get address by id status is 200': (r) => r.status === 200,
      'get address by id returns data': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.data && body.data.address_id;
        } catch (e) {
          return false;
        }
      },
    });

    errorRate.add(!success);
  }

  sleep(1);

  // 12. Update Address
  if (token && contactId && addressId) {
    const payload = JSON.stringify({
      street: `Jl. Updated ${randomString(8)}`,
      city: 'Jakarta Selatan',
      province: 'DKI Jakarta',
      country: 'Indonesia',
      postal_code: '12345',
    });

    const params = {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': token,
      },
    };

    const res = http.put(`${BASE_URL}/address/${contactId}/${addressId}`, payload, params);
    
    const success = check(res, {
      'update address status is 200': (r) => r.status === 200,
      'update address returns updated data': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.message === 'Address updated successfully';
        } catch (e) {
          return false;
        }
      },
    });

    errorRate.add(!success);
  }

  sleep(1);

  // 13. Delete Address
  if (token && contactId && addressId) {
    const params = {
      headers: {
        'Authorization': token,
      },
    };

    const res = http.del(`${BASE_URL}/address/${contactId}/${addressId}`, null, params);
    
    const success = check(res, {
      'delete address status is 200': (r) => r.status === 200,
    });

    errorRate.add(!success);
  }

  sleep(1);

  // 14. Delete Contact
  if (token && contactId) {
    const params = {
      headers: {
        'Authorization': token,
      },
    };

    const res = http.del(`${BASE_URL}/contact/${contactId}`, null, params);
    
    const success = check(res, {
      'delete contact status is 200': (r) => r.status === 200,
    });

    errorRate.add(!success);
  }

  // Random sleep between iterations (1-3 seconds)
  sleep(Math.random() * 2 + 1);
}
