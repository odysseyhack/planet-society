//
//  TransactionPersonalDetailsViewController.swift
//  Permission Hub
//
//  Created by Corné on 13/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class TransactionPluginTableViewCell: UITableViewCell {

    // MARK: - Private properties

    private lazy var stackView: UIStackView = {

        let stackView = UIStackView()
        stackView.translatesAutoresizingMaskIntoConstraints = false

        stackView.axis = .horizontal
        stackView.spacing = 10
        stackView.alignment = .center
        stackView.distribution = .equalSpacing

        return stackView
    }()

    private lazy var digiDButton: UIButton = {

        let button = UIButton()
        button.translatesAutoresizingMaskIntoConstraints = false

        let image = UIImage(named: "digid_button")
        button.setImage(image, for: .normal)

        button.addTarget(
            self,
            action: #selector(digiDButtonTapped),
            for: .touchUpInside)

        return button
    }()

    private lazy var descriptionLabel: UILabel = {

        let label = UILabel()
        label.translatesAutoresizingMaskIntoConstraints = false
        label.font = PHFonts.regular()
        label.textColor = PHColors.greyishBrown
        label.numberOfLines = 0

        label.text = "Use the external DigiD plug-in to fill in your personal information (optional)."

        return label
    }()

    private lazy var separatorView: UIView = {

        let view = UIView()
        view.translatesAutoresizingMaskIntoConstraints = false

        view.backgroundColor = .white

        return view
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

    private func configure() {

        selectionStyle = .none

        addSubview(stackView)
        addSubview(separatorView)

        let margin: CGFloat = 15
        stackView.topAnchor.constraint(equalTo: topAnchor, constant: margin).isActive = true
        stackView.leftAnchor.constraint(equalTo: leftAnchor, constant: margin).isActive = true
        stackView.rightAnchor.constraint(equalTo: rightAnchor, constant: -margin).isActive = true

        configureSeparatorView()

        stackView.addArrangedSubview(digiDButton)
        stackView.addArrangedSubview(descriptionLabel)
    }

    private func configureSeparatorView() {

        let margin: CGFloat = 15
        separatorView.topAnchor.constraint(equalTo: stackView.bottomAnchor, constant: margin).isActive = true
        separatorView.leftAnchor.constraint(equalTo: leftAnchor).isActive = true
        separatorView.rightAnchor.constraint(equalTo: rightAnchor).isActive = true
        separatorView.bottomAnchor.constraint(equalTo: bottomAnchor).isActive = true
        separatorView.heightAnchor.constraint(equalToConstant: 1).isActive = true
    }

    // MARK: - Selectors

    @objc private func digiDButtonTapped(_ sender: UIButton) {
        // ...
    }
}

final class TransactionPersonalDetailsViewController: PHTableViewController {

    // MARK: - Private properties

    private var transaction: TransactionNotification

    // MARK: - Initialization

    init(transaction: TransactionNotification) {
        self.transaction = transaction

        super.init(title: "", items: [
            .notification(
                type: .verification,
                text: "This company is verified"),
            .description(
                date: transaction.date,
                title: "Personal details",
                description: "Please fill out your personal details."),
            .plugin,
            .form(placeholder: "First name"),
            .form(placeholder: "Last name"),
            .form(placeholder: "Date of birth"),
            .form(placeholder: "Address"),
            .form(placeholder: "Email"),
            .form(placeholder: "BSN number")
        ])
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func viewDidLoad() {
        super.viewDidLoad()

        tableView.allowsSelection = false
    }
}
