//
//  TransactionWalletViewController.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights

import UIKit

struct WalletItem {
    let image: UIImage?
    let title: String
    let subtitle: String
    let date: Date
}

final class WalletItemTableViewCell: PHBaseTableViewCell {

    // MARK: - Private properties

    private lazy var horizontalStackView: UIStackView = {

        let stackView = UIStackView()
        stackView.translatesAutoresizingMaskIntoConstraints = false

        stackView.axis = .horizontal
        stackView.spacing = 10
        stackView.alignment = .center

        return stackView
    }()

    private lazy var verticalStackView: UIStackView = {

        let stackView = UIStackView()
        stackView.translatesAutoresizingMaskIntoConstraints = false

        stackView.axis = .vertical
        stackView.spacing = 5
        stackView.alignment = .leading

        return stackView
    }()

    private lazy var itemImageView: UIImageView = {

        let imageView = UIImageView()
        imageView.translatesAutoresizingMaskIntoConstraints = false
        imageView.setContentHuggingPriority(.required, for: .horizontal)

        return imageView
    }()

    private lazy var itemLabel: UILabel = {

        let label = UILabel()
        label.font = PHFonts.regular()
        label.textColor = PHColors.greyishBrown
        label.textAlignment = .left

        return label
    }()

    private let itemSubtitleLabel: UILabel = {

        let label = UILabel()
        label.font = PHFonts.regular()
        label.textColor = PHColors.grey

        return label
    }()

    private lazy var dateLabel: UILabel = {

        let label = UILabel()
        label.font = PHFonts.regular(ofSize: 11)
        label.textColor = PHColors.grey

        label.setContentHuggingPriority(.required, for: .horizontal)

        return label
    }()

    // MARK: - Initialization

    override init(style: UITableViewCell.CellStyle, reuseIdentifier: String?) {
        super.init(style: style, reuseIdentifier: reuseIdentifier)

        configure()
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Configuration

    override func configure() {
        super.configure()

        selectionStyle = .none
        accessoryType = .disclosureIndicator

        addSubview(horizontalStackView)
        separatorView.backgroundColor = .white

        let margin: CGFloat = 15
        horizontalStackView.topAnchor.constraint(equalTo: topAnchor, constant: margin).isActive = true
        horizontalStackView.leftAnchor.constraint(equalTo: leftAnchor, constant: margin).isActive = true
        horizontalStackView.rightAnchor.constraint(equalTo: rightAnchor, constant: -50).isActive = true
        horizontalStackView.bottomAnchor.constraint(equalTo: bottomAnchor, constant: -margin).isActive = true

        horizontalStackView.addArrangedSubview(itemImageView)
        horizontalStackView.addArrangedSubview(verticalStackView)
        verticalStackView.addArrangedSubview(itemLabel)
        verticalStackView.addArrangedSubview(itemSubtitleLabel)
        horizontalStackView.addArrangedSubview(dateLabel)
    }

    func configure(withItem item: WalletItem) {

        itemImageView.image = item.image
        itemLabel.text = item.title
        itemSubtitleLabel.text = item.subtitle
        dateLabel.text = item.date.dateString()
    }
}

final class TransactionWalletViewController: PHTableViewController {

    // MARK: - Initialization

    init() {
        super.init(items: [PHTableViewViewCellType]())
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func viewDidLoad() {
        super.viewDidLoad()

        navigationItem.title = "My Permission Wallet"

        tableView.register(
            WalletItemTableViewCell.self,
            forCellReuseIdentifier: String(describing: WalletItemTableViewCell.self))

        NotificationCenter.default.addObserver(
            self,
            selector: #selector(addItemToWallet),
            name: Notification.Name("add item to wallet"),
            object: nil)
    }

    override func viewDidAppear(_ animated: Bool) {
        super.viewDidAppear(animated)

        // populate view with mock data
        let data = "{\"transactionID\":\"b80db272b05b9ad007c6833dac68b95ca907594946b2da1929d1f8f95d973b5c\",\"item\":[{\"Item\":\"Access to your Personal Details\",\"Fields\":[\"Name\",\"Surname\",\"Date of birth\",\"Email\",\"BSN\"]},{\"Item\":\"Legal identity (passport)\",\"Fields\":[\"Number\",\"Expiration date\",\"Country of issue\"]},{\"Item\":\"Newsletter\",\"Fields\":[\"Email address for marketing purposes\"]},{\"Item\":\"Payment information\",\"Fields\":[\"IBAN number\",\"Bank name\",\"Payment details\"]},{\"Item\":\"Subscription contract 24 months\",\"Fields\":[\"Read and accept the terms\"]}],\"title\":\"Provide permission for completing\",\"description\":\"T-mobile monthly plan(unlimited data), 65 euro, iPhone XR 256GB\",\"verification\":[\"digid.nl\",\"planet-blockchain\",\"kvk\"],\"date\":\"2019-04-13T15:51:57+02:00\",\"requesterName\":\"John Smith\",\"RequesterPublicKey\":\"69093eef7426963f2ef0f68fb73e355b7898ddb04a4fad769a96b41ffc824c1c\",\"analysis\":[\"personal data is GDPR protected data\",\"banking details is sensitive data\"]}".data(using: .utf8)!
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        let transaction = try! decoder.decode(TransactionNotification.self, from: data)
//        presentTransactionFlowViewController(transaction: transaction)
    }

    override func viewWillAppear(_ animated: Bool) {
        super.viewWillAppear(animated)

        pollForNotifications()
    }

    // MARK: - Selectors

    @objc private func addItemToWallet() {

        self.items = [
            .walletItem(item: WalletItem(
                image: UIImage(named: "phone_house"),
                title: "Purchase at Phone House",
                subtitle: "Personal details, Passport, Bank account",
                date: Date()))
        ]

        tableView.reloadData()
    }

    // MARK: - Networking

    @objc private func pollForNotifications() {

        let service = NetworkingService.shared
        try! service.getNotifications { result in

            switch result {
            case .success(let transactionOrNil):
                if let transaction = transactionOrNil {
                    self.presentTransactionFlowViewController(transaction: transaction)
                }
            case .failure(let error):
                print(error)
            }
        }

        // poll endpoint every second
        perform(#selector(pollForNotifications), with: nil, afterDelay: 1)
    }

    // MARK: - Helpers

    private func presentTransactionFlowViewController(transaction: TransactionNotification) {

        let viewController = TransactionFlowViewController(transaction: transaction)
        let navigationController = UINavigationController(rootViewController: viewController)
        present(navigationController, animated: true)
    }
}
