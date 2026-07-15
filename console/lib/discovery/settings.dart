import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'settings.video.dart';
import 'settings.audio.dart';
import 'settings.locate.dart';
import 'api.dart' as api;

class Settings extends StatelessWidget {
  const Settings({super.key});

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return Column(
      mainAxisSize: MainAxisSize.min,
      spacing: defaults.spacing / 2.5,
      children: [
        ds.Heading(Text("discovery")),
        LocateSettings.future(
          api.configuration.get(options: [authn.request(authn.AuthzCache.meta(context))]),
          onChange: (v) => api.configuration.create(v, options: [authn.request(authn.AuthzCache.meta(context))]),
        ),
        ds.Heading(Text("video")),
        VideoSettings(),
        ds.Heading(Text("audio")),
        SettingsAudio(),
      ],
    );
  }
}
