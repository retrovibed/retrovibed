import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import './sku.dart';
import './api.dart' as quotas;
import './typography.dart' as _typography;

class Card extends StatefulWidget {
  final EdgeInsets margin;
  const Card({super.key, this.margin = EdgeInsets.zero});

  @override
  State<Card> createState() => _CardState();
}

class _CardState extends State<Card> {
  Future<quotas.Quota> _storageFuture = Future.value(quotas.Quota());
  Future<quotas.Quota> _bandwidthFuture = Future.value(quotas.Quota());

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  @override
  void initState() {
    super.initState();
    _storageFuture = httpx.withRetry(
      () => quotas.api
          .get(Storage.sku, options: [authn.Authenticated.bearer(context)])
          .then((v) => v.quota..description = Storage.description)
          .catchError((cause) {
            return Storage;
          }, test: httpx.ErrorsTest.err404),
    );

    _bandwidthFuture = httpx.withRetry(
      () => quotas.api
          .get(Bandwidth.sku, options: [authn.Authenticated.bearer(context)])
          .then((v) => v.quota..description = Bandwidth.description)
          .catchError((cause) {
            return Bandwidth;
          }, test: httpx.ErrorsTest.err404),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return ds.Card(
      alignment: Alignment.center,
      margin: widget.margin,
      help: ds.Hint(const Text("check storage usage and limits")),
      Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text("Usage", style: theme.textTheme.titleMedium),
          Column(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              _typography.Typography.future(_storageFuture, defaults: Storage),
              _typography.Typography.future(
                _bandwidthFuture,
                defaults: Bandwidth,
              ),
            ],
          ),
        ],
      ),
    );
  }
}
