import 'package:flutter/material.dart';
import 'package:fractal/designkit.dart' as ds;
import 'package:fractal/rss.dart' as rss;
import 'package:fractal/torrents.dart' as torrents;

class Display extends StatelessWidget {
  const Display({super.key});

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return ListView(
      children: [
        ds.Accordion(
          description: Text("RSS"),
          content: Container(
            padding: defaults.padding,
            child: rss.ListSearchable(),
          ),
        ),
        ds.Accordion(
          disabled: Text("coming soon"),
          description: Text("storage"),
          content: Container(),
        ),
        ds.Accordion(
          disabled: Text("coming soon"),
          description: Row(children: [Text("torrents")]),
          content: Column(children: [torrents.SettingsLeech()]),
        ),
        ds.Accordion(
          disabled: Text("coming soon"),
          description: Row(children: [Text("VPN (wireguard)")]),
          content: Container(),
        ),
        ds.Accordion(
          disabled: Text("coming soon - opt in premium features"),
          description: Text("billing"),
          content: Container(),
        ),
      ],
    );
  }
}
