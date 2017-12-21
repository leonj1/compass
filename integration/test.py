#!/usr/bin/python

import requests
import unittest
import json
import uuid

host = "http://localhost:3244"


def my_random_string(string_length=10):
    """Returns a random string of length string_length."""
    random = str(uuid.uuid4()) # Convert UUID format to a Python string.
    random = random.upper() # Make all characters uppercase.
    random = random.replace("-","") # Remove the UUID '-'.
    return random[0:string_length] # Return the random string.


class MyTest(unittest.TestCase):
    # @unittest.skip("testing skipping")
    def test_health_check(self):
        r = requests.get("{}/public/health".format(host))
        self.assertEqual(r.status_code, 200)

    # @unittest.skip("testing skipping")
    def test_add_confession(self):
        r = requests.delete("{}/clusters/dev01".format(host))

        t = "2017-01-02 15:04:05"
        payload = {
           "name": "dev01",
           "status": "live",
           "personality": "dev",
           "events": "this is one line\nThis is another",
           "crds": {
             "pyjob": {
               "name": "pyjob",
               "version": "1.0"
             },
             "sparkjob": {
               "name": "sparkjob",
               "version": "2.0"
             }
           },
           "nodes": {
             "node1": {
               "name": "node1",
               "version": "1.8.0"
             },
             "node2": {
               "name": "node2",
               "version": "1.8.0"
             }
           },
           "namespace": {
              "default": {
                 "name": "default",
                 "pod_count": 3,
                 "crds": {
                    "pyjob": 2,
                    "sparkjob": 3
                 }
              }
           }
         }

        r = requests.post("{}/clusters".format(host), json=payload)
        self.assertEqual(r.status_code, 201)
        r = requests.get("{}/clusters/dev01".format(host))
        self.assertEqual(r.status_code, 200)
        expected = '{"name":"dev01","status":"live","personality":"dev","crds":{"pyjob":{"name":"pyjob","version":"1.0"},"sparkjob":{"name":"sparkjob","version":"2.0"}},"nodes":{"node1":{"name":"node1","version":"1.8.0"},"node2":{"name":"node2","version":"1.8.0"}},"namespace":{"default":{"name":"default","pod_count":3,"crds":{"pyjob":2,"sparkjob":3}}},"events":"this is one line\\nThis is another"}'
        self.assertEqual(r.content, expected)
        # Add Crd
        payload = {
           "name": "tensorflow",
           "version": "1.9.0"
        }
        r = requests.post("{}/clusters/dev01/crds".format(host), json=payload)
        self.assertEqual(r.status_code, 201)
        expected = '{"message":"created"}'
        self.assertEqual(r.content, expected)
        r = requests.get("{}/clusters/dev01".format(host))
        self.assertEqual(r.status_code, 200)
        expected = '{"name":"dev01","status":"live","personality":"dev","crds":{"pyjob":{"name":"pyjob","version":"1.0"},"sparkjob":{"name":"sparkjob","version":"2.0"},"tensorflow":{"name":"tensorflow","version":"1.9.0"}},"nodes":{"node1":{"name":"node1","version":"1.8.0"},"node2":{"name":"node2","version":"1.8.0"}},"namespace":{"default":{"name":"default","pod_count":3,"crds":{"pyjob":2,"sparkjob":3}}},"events":"this is one line\\nThis is another"}'
        self.assertEqual(r.content, expected)
        # Add Node
        payload = {
           "name": "hostA",
           "version": "1.9.0"
        }
        r = requests.post("{}/clusters/dev01/nodes".format(host), json=payload)
        self.assertEqual(r.status_code, 201)
        expected = '{"message":"created"}'
        self.assertEqual(r.content, expected)
        r = requests.get("{}/clusters/dev01".format(host))
        self.assertEqual(r.status_code, 200)
        expected = '{"name":"dev01","status":"live","personality":"dev","crds":{"pyjob":{"name":"pyjob","version":"1.0"},"sparkjob":{"name":"sparkjob","version":"2.0"},"tensorflow":{"name":"tensorflow","version":"1.9.0"}},"nodes":{"hostA":{"name":"hostA","version":"1.9.0"},"node1":{"name":"node1","version":"1.8.0"},"node2":{"name":"node2","version":"1.8.0"}},"namespace":{"default":{"name":"default","pod_count":3,"crds":{"pyjob":2,"sparkjob":3}}},"events":"this is one line\\nThis is another"}'
        self.assertEqual(r.content, expected)
        # Add Namespace
        payload = {
           "name": "platform",
           "pods": "9",
           "crds": {
              "crd1": 1,
              "crd2": 2,
              "crd3": 3
           }
        }
        r = requests.post("{}/clusters/dev01/namespaces".format(host), json=payload)
        self.assertEqual(r.status_code, 201)
        expected = '{"message":"created"}'
        self.assertEqual(r.content, expected)
        r = requests.get("{}/clusters/dev01".format(host))
        self.assertEqual(r.status_code, 200)
        expected = '{"name":"dev01","status":"live","personality":"dev","crds":{"pyjob":{"name":"pyjob","version":"1.0"},"sparkjob":{"name":"sparkjob","version":"2.0"},"tensorflow":{"name":"tensorflow","version":"1.9.0"}},"nodes":{"hostA":{"name":"hostA","version":"1.9.0"},"node1":{"name":"node1","version":"1.8.0"},"node2":{"name":"node2","version":"1.8.0"}},"namespace":{"default":{"name":"default","pod_count":3,"crds":{"pyjob":2,"sparkjob":3}},"platform":{"name":"platform","crds":{"crd1":1,"crd2":2,"crd3":3}}},"events":"this is one line\\nThis is another"}'
        self.assertEqual(r.content, expected)


if __name__ == '__main__':
    unittest.main()

