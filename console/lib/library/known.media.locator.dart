import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'known.media.card.dart';
import './api.dart' as api;

const String _locateDisclaimerCacheId = 'discovery.p2p';

const String _locateDisclaimerText =
    'Retrovibed operates on an entirely peer-to-peer (p2p) system.\n\n'
    'The p2p system allows journalists, musicians, and other content creators to directly '
    'interact with users and be discovered, preventing censorship by third parties '
    '(including retrovibed).\n\n'
    'As a result, it can locate media that has been published by third parties who do not '
    'have distribution rights. Retrovibed takes no responsibility for such content.\n\n'
    'P2P: operate in fully p2p mode, allowing discovery of publicly available content. '
    'You take responsibility for your activities and for obeying all legal laws of your region.\n\n'
    'Listed Only: only allows content published directly to the retrovibed index to work.';

enum _DiscoveryChoice { nevermind, p2p, listedOnly }

class _LocateDisclaimerPrompt extends StatelessWidget {
  final void Function(_DiscoveryChoice) onChoice;
  const _LocateDisclaimerPrompt({required this.onChoice});

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return Container(
      padding: defaults.padding,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        spacing: defaults.spacing,
        children: [
          Expanded(
            child: SingleChildScrollView(child: Text(_locateDisclaimerText)),
          ),
          Row(
            spacing: defaults.spacing,
            mainAxisSize: MainAxisSize.min,
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              TextButton(onPressed: () => onChoice(_DiscoveryChoice.nevermind), child: Text('Nevermind')),
              TextButton(onPressed: () => onChoice(_DiscoveryChoice.p2p), child: Text('P2P')),
              TextButton(onPressed: () => onChoice(_DiscoveryChoice.listedOnly), child: Text('Listed Only')),
            ],
          ),
        ],
      ),
    );
  }
}

class KnownMediaLocator extends StatefulWidget {
  final api.Known current;
  final Future<api.LocateCreateResponse> Function(api.Locate req, {List<httpx.Option> options}) locate;
  final bool Function(String)? disclaimer;
  final void Function(String)? acknowledge;
  final IconData icon;
  final Widget help;
  final Widget? trailing;

  const KnownMediaLocator(
    this.current, {
    super.key,
    this.locate = api.locate.create,
    this.icon = Icons.search,
    this.help = ds.HelpScope.None,
    this.disclaimer,
    this.acknowledge,
    this.trailing,
  });

  @override
  State<StatefulWidget> createState() => _KnownMediaLocator();
}

class _KnownMediaLocator extends State<KnownMediaLocator> {
  bool _loading = false;
  Widget _cause = ds.Error.zero;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  void _onDisclaimerChoice(_DiscoveryChoice choice) {
    if (choice == _DiscoveryChoice.nevermind) return;
    (widget.acknowledge ?? ds.Disclaimer.acknowledge)(_locateDisclaimerCacheId);
    setState(() {});
  }

  void _onTap() async {
    setState(() {
      _loading = true;
      _cause = ds.Error.zero;
    });

    widget
        .locate(
          api.Locate.create()..knownMediaId = widget.current.id,
          options: [authn.request(authn.AuthzCache.meta(context))],
        )
        .then((_) {
          setState(() {
            _loading = false;
          });
        })
        .catchError((e) {
          setState(() {
            _loading = false;
            _cause = ds.Error.unknown(e, onTap: reseterr);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    return ds.Loading(
      loading: _loading,
      cause: _cause,
      ds.Disclaimer(
        KnownMediaCard(
          widget.current,
          icon: widget.icon,
          help: widget.help,
          onTap: _loading ? null : _onTap,
          trailing: widget.trailing,
        ),
        overlay: _LocateDisclaimerPrompt(onChoice: _onDisclaimerChoice),
        cacheid: _locateDisclaimerCacheId,
        cached: widget.disclaimer ?? ds.Disclaimer.disclaimerpath,
      ),
    );
  }
}
