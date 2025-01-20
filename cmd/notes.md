                                     ##### `<b>Client </b>`
                                            -Input form
                                            -Websocket connection
                                                      |
                                                      | send message via ws
                                                      |
                                            <b>Websocket server</b>
                                            -Accept websocket connection
                                            -Handle client messages
                                            -Broadcast to clients
                                            -Publish to redis
                                                      |
                                                      | put message in db
                                                      |
                                            <b>Redis database</b>
                                            -Store messages
                                            -Publish messages to "chat:messages" channel
                                                      |
                                                      | broadcast to websocket server
                                                      |
                                            <b>Websocket server</b>
                                            -Handle websocket connection.
                                            -Broadcasdt to clients.
                                                      |
                                                      | websocket message
                                                      |
                                            <b>Client</b>
                                            -Display messages