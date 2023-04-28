import { createServer } from "miragejs"

export function makeServer({ environment = "development" } = {}) {
    return createServer({
        environment,
        routes(){
            this.post("/api/login", (schema, request) => {
                const body = JSON.parse(request.requestBody)
                const username = body.username
                const password = body.password
                if('dev' === username && "dev" === password){
                  return {username: body.username, token: 'tk-0000-0000'}
                }
                return {}
            })
            this.get('/api/wallets', (schema, request) => {
                return [
                    {name:'do-vv',   address:'f14ognplttj7lyag7wug77auazvs2t6wg7xlvepay', balance:270.29, nodes:1, daily:10},
                    {name:'do-xue',  address:'f1tkr2tz5nojynvu5bs54uzgffluklo65ba3gruba', balance:259.79, nodes:1, daily:11},
                    {name:'id-west', address:'f1xbawvqfdgu5ptfejxgqx7vbqeypjrgee7i3w4qi', balance:290.29, nodes:2, daily:13},
                    {name:'id-eno',  address:'f15qiacu22v74vbrvrgcevtt7wezqa5wrxbsjtg6y', balance:260.09, nodes:4, daily:14},
                ]
            }),
            this.post('/api/wallets', (schema, request) => {
                const body = JSON.parse(request.requestBody)
                const name = body.name
                const address = body.address
                console.log(name, address)
                return "ok"
            }),
            this.get('/api/summary', (schema, request) => {
                return {
                  "earnings": "2300",
                  "nodes": 11,
                  "inactive": 1,
                  "dailys": [34, 45, 20, 69, 33, 98, 23, 33, 10, 33, 70, 100],
                  "time":"2023-04-09 15:03:51",
                }
            }),
            this.get('/api/nodes', (schema, request) => {
                return [
                    {id:'423667f8', ip:'44.234.105.210', isp:'Oracle', location:'India', cpu:6, ram:24, disk:200, bandwidth:100, traffic: 5},
                    {id:'f5d3fc3c', ip:'44.234.105.210', isp:'Melbikomas UAB', location:'Singapore (SG)', cpu:4, ram:24, disk:300, bandwidth:4000, traffic: 0},
                ]
            })
        }
    })
}