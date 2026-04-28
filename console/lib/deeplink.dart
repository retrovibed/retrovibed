import 'dart:async';
import 'package:app_links/app_links.dart';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/billing/api.dart' as billing;
import 'package:retrovibed/community/api.dart' as community;
import 'package:retrovibed/community/community.detail.dart';
import 'package:retrovibed/community/link.content.dart';

typedef SearchCommunity =
    Future<community.CommunitySearchResponse> Function(
      community.CommunitySearchRequest req, {
      List<httpx.Option> options,
    });
typedef ConsumeAttribution =
    Future<billing.AttributionConsumeResponse> Function(String token, {List<httpx.Option> options});
typedef SubscribeAction =
    Future<void> Function(BuildContext context, community.Community community, String attribution);

Stream<Uri> _defaultUriStream() => AppLinks().uriLinkStream;
Future<Uri?> _defaultInitialUri() => AppLinks().getInitialLink();

class DeepLink extends StatefulWidget {
  final Widget child;
  final SearchCommunity search;
  final ConsumeAttribution consumeAttribution;
  final SubscribeAction subscribe;
  final Stream<Uri> Function() uriStream;
  final Future<Uri?> Function() initialUri;

  const DeepLink(
    this.child, {
    super.key,
    this.search = community.API.search,
    this.consumeAttribution = billing.consumeAttribution,
    this.subscribe = handleSubscribeAction,
    this.uriStream = _defaultUriStream,
    this.initialUri = _defaultInitialUri,
  });

  @override
  State<DeepLink> createState() => _DeepLinkState();
}

class _DeepLinkState extends State<DeepLink> {
  StreamSubscription<Uri>? _sub;
  Widget _overlay = const SizedBox();

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  @override
  void initState() {
    super.initState();
    _sub = widget.uriStream().listen(_handleUri);
    widget.initialUri().then((uri) {
      if (uri != null) _handleUri(uri);
    }).ignore();
  }

  @override
  void dispose() {
    _sub?.cancel();
    super.dispose();
  }

  void _dismiss() {
    setState(() => _overlay = const SizedBox());
  }

  void _handleUri(Uri uri) {
    final host = uri.host;
    if (host == 'invite.retrovibe.space') {
      _handleInvite(uri);
      return;
    }
    if (host.endsWith('.community.retrovibe.space')) {
      _handleCommunity(uri);
      return;
    }
  }

  void _handleInvite(Uri uri) {
    final attribution = uri.queryParameters['a'] ?? '';
    if (attribution.isEmpty) return;

    widget
        .consumeAttribution(
          attribution,
          options: [authn.Authenticated.bearer(context)],
        )
        .ignore();
  }

  void _handleCommunity(Uri uri) {
    final domain = uri.host.split('.').first;
    if (domain.isEmpty) return;

    httpx
        .withRetry(
          () => widget.search(
            community.CommunitySearchRequest(query: domain),
            options: [authn.AuthzCache.bearer(context)],
          ),
        )
        .then((response) {
          if (response.items.isEmpty) {
            setState(
              () =>
                  _overlay = ds.Masked(
                    alignment: Alignment.center,
                    ds.Error.unknown('community not found', onTap: _dismiss),
                  ),
            );
            return;
          }

          final c = response.items.first;
          setState(
            () =>
                _overlay = ds.Masked(
                  alignment: Alignment.center,
                  ds.Confirmation.yesNo(
                    content: Column(
                      children: [CommunityDetail(community: c)],
                    ),
                    onConfirm: () {
                      widget.subscribe(context, c, '').catchError((e, s) {
                        setState(
                          () => _overlay = ds.Error.unknown(e, onTap: _dismiss),
                        );
                        return Future.error(e);
                      }).ignore();
                      _dismiss();
                    },
                    onCancel: _dismiss,
                  ),
                ),
          );
        })
        .catchError((e) {
          setState(
            () =>
                _overlay = ds.Masked(
                  alignment: Alignment.center,
                  ds.Error.unknown(e, onTap: _dismiss),
                ),
          );
        });
  }

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        widget.child,
        _overlay,
      ],
    );
  }
}
