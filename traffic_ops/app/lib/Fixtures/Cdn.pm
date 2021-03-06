package Fixtures::Cdn;
#
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
use Moose;
extends 'DBIx::Class::EasyFixture';
use namespace::autoclean;
use Digest::SHA1 qw(sha1_hex);

my %definition_for = (
	cdn1_cdn_name => {
		new   => 'Cdn',
		using => {
			id          => 100,
			name        => 'cdn1',
			domain_name => 'cdn1.kabletown.net',
		},
	},
	cdn2_cdn_name => {
		new   => 'Cdn',
		using => {
			id          => 200,
			name        => 'cdn2',
			domain_name => 'cdn2.kabletown.net',
		},
	},
	cdn3_cdn_name => {
		new   => 'Cdn',
		using => {
			id          => 300,
			name        => 'cdn3',
			domain_name => 'cdn3.kabletown.net',
		},
	},
);

sub get_definition {
	my ( $self, $name ) = @_;
	return $definition_for{$name};
}

sub all_fixture_names {
	# sort by db id to guarantee insert order
	return (sort { $definition_for{$a}{using}{id} cmp $definition_for{$b}{using}{id} } keys %definition_for);
}

__PACKAGE__->meta->make_immutable;

1;
