import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/meta.dart' as meta;

class Display extends StatefulWidget {
  const Display({super.key});

  @override
  State<Display> createState() => _DeepLinkState();
}

class _DeepLinkState extends State<Display> {
  ValueNotifier<meta.Daemon> _library = ValueNotifier(meta.Daemon());

  void refresh() {
    setState(() {}); // force rebuild
  }

  @override
  void initState() {
    super.initState();
    _library = meta.EndpointAuto.of(context)?.changed ?? _library;
    _library.addListener(refresh);
  }

  @override
  void dispose() {
    super.dispose();
    _library.removeListener(refresh);
  }

  @override
  Widget build(BuildContext context) {
    return ds.Container(
      meta.DaemonDropdown(
        library: _library,
        trailing: [],
      ),
    );
  }
}
