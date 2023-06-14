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
            this.get('/api/wallets', () => {
                return [
                    {name:'do-vv',   address:'f14ognplttj7lyag7wug77auazvs2t6wg7xlvepay', balance:270.29, nodes:[10,2], daily:10, group:'us'},
                    {name:'do-xue',  address:'f1tkr2tz5nojynvu5bs54uzgffluklo65ba3gruba', balance:259.79, nodes:[10,1], daily:11, group:'asia'},
                    {name:'id-west', address:'f1xbawvqfdgu5ptfejxgqx7vbqeypjrgee7i3w4qi', balance:290.29, nodes:[10,0], daily:13, group:'asia'},
                    {name:'id-eno',  address:'f15qiacu22v74vbrvrgcevtt7wezqa5wrxbsjtg6y', balance:260.09, nodes:[100,0], daily:14, group:'us'},
                ]
            })
            this.post('/api/wallets', (schema, request) => {
                const body = JSON.parse(request.requestBody)
                const name = body.name
                const address = body.address
                const group = body.group
                console.log('server:', name, address, group)
                return "ok"
            })
            this.get('/api/summary', () => {
                return {
                  "cost": "2300",
                  "roi": 250,
                  "earnings": "2300",
                  "nodes": [100,10],
                  "dailys": [34, 45, 20, 69, 33, 98, 23, 33, 10, 33, 70, 100],
                  "time":"2023-04-09 15:03:51",
                  "groups": [{name:'-', balance:28.947248},{name:'linode', balance:6.405105},{name:'asia', balance:1.30295024}]
                }
            })
            this.get('/api/nodes', () => {
                return [
                    {id:'423667f8', name:'aws-id-1',state:'active','type':883, ip:'44.214.105.1', isp:'Oracle', location:'India', cpu:6, ram:24, disk:'228G/249G (92%)', bandwidth:10, traffic: '0/20TB'},
                    {id:'423667f8', name:'aws-id-2',state:'draining','type':883, ip:'44.234.103.2', isp:'Oracle', location:'India', cpu:6, ram:24, disk:'228G/249G (92%)', bandwidth:10, traffic: '160.99 MiB/20TB'},
                    {id:'423667f8', name:'aws-id-3',state:'down','type':883, ip:'44.234.106.3', isp:'Oracle', location:'India', cpu:6, ram:24, disk:'228G/249G (92%)', bandwidth:10, traffic: '0/20TB', price:45},
                    {id:'423667f8', name:'aws-id-4',state:'','type':883, ip:'44.234.108.4', isp:'Oracle', location:'India', cpu:6, ram:24, disk:'228G/249G (92%)', bandwidth:10, traffic: '0/20TB', renew:'2013-07-01'},
                ]
            })
            this.post('/api/nodes/pool', () => {
                return "ok"
            })
        }
    })
}