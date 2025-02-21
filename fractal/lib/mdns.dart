import 'package:flutter/material.dart';
import 'package:multicast_dns/multicast_dns.dart';
import 'package:fractal/designkit.dart' as ds;
import './httpx.dart' as httpx;

class MDNSDiscovery extends StatefulWidget {
  final Widget child;
  final UniqueKey id = UniqueKey();

  MDNSDiscovery(this.child, {super.key});

  static _MDNSDiscovery? of(BuildContext context) {
    return context.findAncestorStateOfType<_MDNSDiscovery>();
  }

  @override
  State<StatefulWidget> createState() => _MDNSDiscovery();
}

class _MDNSDiscovery extends State<MDNSDiscovery> {
  static const String ServiceName = "_shallows._udp.local";
  // Key _refresh = UniqueKey();
  final MDnsClient _client = MDnsClient();
  var _loading = true;

  @override
  void initState() {
    super.initState();
    _client
        .start()
        .then((_) {
          print("mDNS started");
          _client
              .lookup<PtrResourceRecord>(
                ResourceRecordQuery.service(ServiceName),
              )
              // .lookup<PtrResourceRecord>(
              //   ResourceRecordQuery.addressIPv4(ServiceName),
              // )
              // .lookup<PtrResourceRecord>(
              //   ResourceRecordQuery.addressIPv6(ServiceName),
              // )
              .listen((event) {
                httpx.set(event.domainName);
                print("waka 0: ${event}");
              });
          _client
              .lookup<PtrResourceRecord>(
                ResourceRecordQuery.serverPointer(ServiceName),
              )
              .listen((ptr) {
                _client
                    .lookup<SrvResourceRecord>(
                      ResourceRecordQuery.service(ptr.domainName),
                    )
                    .listen((srv) {
                      print("waka 1: ${srv}");
                      httpx.set("${srv.target}:${srv.port}");
                      setState(() {
                        _loading = false;
                      });
                    });
              });
        })
        .catchError((cause) {
          print("mDNS failed: $cause");
        });
  }

  @override
  void dispose() {
    super.dispose();
    _client.stop();
  }

  @override
  Widget build(BuildContext context) {
    return ds.Loading(loading: _loading, child: widget.child);
    // return Container(key: _refresh, child: widget.child);
  }
}
