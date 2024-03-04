require 'digest/md5'

class Key < ActiveRecord::Base
  include Gitlab::CurrentSettings
  include Sortable

  belongs_to :user

  before_validation :generate_fingerprint

  validates :title,
    presence: true,
    length: { maximum: 255 }

  validates :key,
    presence: true,
    length: { maximum: 5000 },
    format: { with: /\A(ssh|ecdsa)-.*\Z/ }

  validates :fingerprint,
    uniqueness: true,
    presence: { message: 'cannot be generated' }

  validate :key_meets_restrictions

  delegate :name, :email, to: :user, prefix: true

  after_commit :add_to_shell, on: :create
  after_create :post_create_hook
  after_create :refresh_user_cache
  after_commit :remove_from_shell, on: :destroy
  after_destroy :post_destroy_hook
  after_destroy :refresh_user_cache

  def key=(value)
    value&.delete!("\n\r")
    value.strip! unless value.blank