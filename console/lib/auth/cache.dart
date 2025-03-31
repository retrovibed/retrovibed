import 'package:flutter/material.dart';
import 'package:console/designkit.dart' as ds;
import 'package:console/meta.dart' as _meta;
import 'package:console/authz.dart' as authz;

class AuthzCache extends StatefulWidget {
  final Widget child;

  const AuthzCache(this.child, {Key? key}) : super(key: key);
  static _AuthzCache? of(BuildContext context) {
    return context.findAncestorStateOfType<_AuthzCache>();
  }

  @override
  State<AuthzCache> createState() => _AuthzCache();
}

class _AuthzCache extends State<AuthzCache> {
  bool _loading = true;
  authz.Cached<_meta.Token> meta = authz.Cached(
    authz.Bearer(_meta.Token(), ""),
    authz.refresh(
      (c) => _meta.authz(c).then((v) {
        return authz.Bearer(v.token, v.bearer);
      }),
      (c, ts) =>
          DateTime.fromMillisecondsSinceEpoch(c.expires.toInt()).isBefore(ts),
    ),
  );

  @override
  void initState() {
    super.initState();
    Future.wait([meta.token()]).whenComplete(() {
      setState(() {
        _loading = false;
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    return ds.Loading(loading: _loading, child: widget.child);
  }
}
