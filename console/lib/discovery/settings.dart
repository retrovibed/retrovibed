import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'settings.video.dart';
import 'settings.audio.dart';

class Settings extends StatelessWidget {
  const Settings({super.key});
  // static api.StorageSettingsResponse zero = api.StorageSettingsResponse();
  // final api.StorageSettingsResponse current;
  // final Future<api.StorageSettingsResponse> Function(
  //   api.StorageSettingsResponse,
  // )? onChange;
  // const MinimalSettings(this.current, {super.key, this.onChange});

  // static FutureBuilder<api.StorageSettingsResponse> future(
  //   Future<api.StorageSettingsResponse> pending, {
  //   Future<api.StorageSettingsResponse> Function(api.StorageSettingsResponse)?
  //   onChange,
  // }) {
  //   return ds.future(MinimalSettings.zero, pending, (v) {
  //     return MinimalSettings(v, key: UniqueKey(), onChange: onChange);
  //   });
  // }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return Column(
      mainAxisSize: MainAxisSize.min,
      spacing: defaults.spacing / 2.5,
      children: [
        ds.Heading(Text("video")),
        VideoSettings(),
        ds.Heading(Text("audio")),
        SettingsAudio(),
      ],
    );
  }
}
