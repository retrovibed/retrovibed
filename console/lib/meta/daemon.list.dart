import 'package:flutter/material.dart';
import 'package:console/designkit.dart' as ds;
import './api.dart' as api;

class DaemonList extends StatefulWidget {
  final void Function(api.Daemon v)? onTap;
  final Future<api.DaemonSearchResponse> Function(api.DaemonSearchRequest)
  search;

  const DaemonList({super.key, this.search = api.daemons.search, this.onTap});

  @override
  State<StatefulWidget> createState() => _DaemonList();
}

class _DaemonList extends State<DaemonList> {
  bool _loading = true;
  ds.Error? _cause = null;
  api.DaemonSearchResponse _res = api.daemons.response();

  void refresh(api.DaemonSearchRequest req) {
    widget
        .search(req)
        .then((v) {
          setState(() {
            _res = v;
            _loading = false;
          });
        })
        .catchError((e) {
          setState(() {
            _cause = ds.Error.unknown(e);
            _loading = false;
          });
        });
  }

  @override
  void initState() {
    super.initState();
    refresh(_res.next);
  }

  @override
  Widget build(BuildContext context) {
    return ds.Table(
      loading: _loading,
      cause: _cause,
      children: _res.items,
      ds.Table.inline<api.Daemon>(
        (v) => _RowDisplay(
          current: v,
          onTap: widget.onTap == null ? null : () => widget.onTap!(v),
        ),
      ),
    );
  }
}

void notimplemented(String s) {
  return print(s);
}

void rowtapdefault() => notimplemented("row tap not implemented");

class _RowDisplay extends StatelessWidget {
  final api.Daemon current;
  final void Function()? onTap;
  const _RowDisplay({required this.current, this.onTap = rowtapdefault});

  @override
  Widget build(BuildContext context) {
    final themex = ds.Defaults.of(context);
    return Container(
      padding: themex.padding,
      child: InkWell(
        onTap: onTap,
        child: Row(
          spacing: themex.spacing!,
          children: [
            Expanded(
              child: Text(current.description, overflow: TextOverflow.ellipsis),
            ),
          ],
        ),
      ),
    );
  }
}
