import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/retrovibed.dart' as retro;

class MeteredToggle extends StatefulWidget {
  const MeteredToggle({super.key});

  @override
  State<MeteredToggle> createState() => _MeteredToggleState();
}

class _MeteredToggleState extends State<MeteredToggle> {
  bool _metered = retro.metered();
  @override
  Widget build(BuildContext context) {
    return ds.LoadingIconButton(
      toggled: _metered,
      icon: Icon(Icons.network_check, color: _metered ? Colors.green : null),
      onPressed: () async {
        final upd = retro.set_metered(!retro.metered());
        setState(() {
          _metered = upd;
        });
      },
    );
  }
}
