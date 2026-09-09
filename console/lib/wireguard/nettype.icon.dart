import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart' as api;

class NettypeIcon extends StatefulWidget {
  final api.Wireguard current;

  final Future<void> Function(api.WireguardNettype v) onTap;
  const NettypeIcon(
    this.current, {
    super.key,
    this.onTap = ds.fnAsyncNoopOnChange,
  });

  @override
  State<NettypeIcon> createState() => _NettypeIconState();
}

class _NettypeIconState extends State<NettypeIcon> {
  late api.WireguardNettype _nettype;

  @override
  void initState() {
    super.initState();
    _nettype = widget.current.nettype;
  }

  @override
  void didUpdateWidget(covariant NettypeIcon oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.current.nettype != widget.current.nettype) {
      _nettype = widget.current.nettype;
    }
  }

  Future<void> _select(api.WireguardNettype v) {
    setState(() {
      _nettype = v;
    });
    return widget.onTap(v);
  }

  Widget _icon(api.WireguardNettype v) {
    switch (v) {
      case api.WireguardNettype.DISTRIBUTION:
        return const Icon(Icons.hub);
      case api.WireguardNettype.SOCIAL:
        return const Icon(Icons.publish);
      default:
        return ds.Empty;
    }
  }

  @override
  Widget build(BuildContext context) {
    return PopupMenuButton<String>(
      position: PopupMenuPosition.under,
      color: Theme.of(context).colorScheme.surface,
      surfaceTintColor: Theme.of(context).colorScheme.surface,
      icon: _icon(_nettype),
      itemBuilder: (context) => [
        PopupMenuItem<String>(
          onTap: () => _select(api.WireguardNettype.UNSPECIFIED),
          child: ds.LoadingListTile(
            leading: ds.Empty,
            title: const Text("Unspecified"),
          ),
        ),
        PopupMenuItem<String>(
          onTap: () => _select(api.WireguardNettype.DISTRIBUTION),
          child: ds.LoadingListTile(
            leading: const Icon(Icons.hub),
            title: const Text("Distribution"),
          ),
        ),
        PopupMenuItem<String>(
          onTap: () => _select(api.WireguardNettype.SOCIAL),
          child: ds.LoadingListTile(
            leading: const Icon(Icons.publish),
            title: const Text("Social"),
          ),
        ),
      ],
    );
  }
}
