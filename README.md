# k6-plugin-template

Plugin example script:

```javascript
import { func } from 'k6-plugin/name';  // import plugin

export default function () {
    func(1); // func with count input variable
}
```

Result output:

```bash
$ sudo ./k6 run --vus 50 --duration 30s --plugin=name.so test.js

          /\      |‾‾|  /‾‾/  /‾/
     /\  /  \     |  |_/  /  / /
    /  \/    \    |      |  /  ‾‾\  
   /          \   |  |‾\  \ | (_) |
  / __________ \  |__|  \__\ \___/ .io

  execution: local
    plugins: <NAME>
     output: -
     script: test.js

    duration: 30s, iterations: -
         vus: 50,

  duration: 30s, iterations: -
         vus: 2,  

  execution: local
     script: test.js
     output: -

  scenarios: (100.00%) 1 executors, 2 max VUs, 1m0s max duration (incl. graceful stop):
           * default: 2 looping VUs for 30s (gracefulStop: 30s)


running (0m30.1s), 0/2 VUs, 560 complete and 0 interrupted iterations
default ✓ [======================================] 2 VUs  30s


    data_received...........: 0 B     0 B/s
    data_sent...............: 0 B     0 B/s
    iteration_duration......: avg=107.29ms min=9.3ms med=110.52ms max=144.55ms p(90)=115.5ms p(95)=119.26ms
    iterations..............: 560     18.594848/s
    vus.....................: 2       min=2 max=2
    vus_max.................: 2       min=2 max=2
```
