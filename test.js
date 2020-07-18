/*

This is a k6 test script that imports the k6-icmp-plugin.

*/

import { check } from 'k6';
import { ping } from 'k6-plugin/icmp';  // import icmp plugin

export default function () {
    const hostname = "google.com";
    const count = 1;
    const interval = 1;
    const timeout = 1;
    const size = 64;
    const error = ping(hostname, count, interval, timeout, size);

    check(error, {
        "ping successful": err => err == undefined
    });
}
