import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/quotas.dart' as quotas;

import './api.dart' as api;
import './local.storage.dart';
import './archive.storage.dart';

class MinimalSettings extends StatelessWidget {
  static api.StorageSettingsResponse zero = api.StorageSettingsResponse();
  final api.StorageSettingsResponse current;
  final Future<api.StorageSettingsResponse> Function(
    api.StorageSettingsResponse,
  )?
  onChange;
  final httpx.Endpoint<String, quotas.QuotaFindResponse> archiveQuota;

  const MinimalSettings(
    this.current, {
    super.key,
    this.onChange,
    this.archiveQuota = quotas.api.get,
  });

  static FutureBuilder<api.StorageSettingsResponse> future(
    Future<api.StorageSettingsResponse> pending, {
    Future<api.StorageSettingsResponse> Function(api.StorageSettingsResponse)? onChange,
    httpx.Endpoint<String, quotas.QuotaFindResponse> archiveQuota = quotas.api.get,
  }) {
    return ds.future(MinimalSettings.zero, pending, (snapshot) {
      return ds.ErrorScreen(
        MinimalSettings(
          snapshot.data ?? MinimalSettings.zero,
          key: ValueKey(snapshot.data.hashCode),
          onChange: onChange,
          archiveQuota: archiveQuota,
        ),
        cause: snapshot.hasError ? ds.Error.unknown(snapshot.error!) : ds.Error.zero,
      );
    });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return Column(
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.stretch,
      spacing: defaults.spacing / 2.5,
      children: [
        ds.Container(
          padding: defaults.padding,
          margin: EdgeInsets.zero,
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: const [Text("library"), Text("cache storage")],
          ),
        ),
        LocalStorageSettings(
          current.local,
          onChange: (l) {
            final _upd = current.deepCopy()..local = l;
            return onChange?.call(_upd).then((v) => v.local) ?? Future.value(current.local);
          },
        ),
        ds.Container(
          padding: defaults.padding,
          margin: EdgeInsets.zero,
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: const [Text("archive"), Text("cloud storage")],
          ),
        ),
        ArchiveStorage(quota: archiveQuota),
      ],
    );
  }
}
